package importer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-dd-dataset-importer/content"
	"github.com/ONSdigital/dp-dd-dataset-importer/wda"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"log"
	"net/http"
	"path"
	"strings"
)

func ImportDataset(datasetSource string, forceDownload bool, indexerUrl string) {

	fmt.Println("Importing dataset: " + datasetSource)
	datasetPath := DownloadDataset(datasetSource, forceDownload)
	dataset := MapDataset(datasetPath)

	// save
	fileName := UrlToFilename(datasetSource)
	outputFilePath := path.Join(OutputDir, fileName)
	content.SaveObjectJson(dataset, outputFilePath)

	// index
	fmt.Println(len(indexerUrl))
	fmt.Println(indexerUrl)
	if len(indexerUrl) > 0 {

		fmt.Println("Sending document to indexer")
		document := &model.Document{
			Body: dataset,
			Type: "dataset",
		}
		jsonBytes, err := json.Marshal(document)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post(indexerUrl, "application/json", bytes.NewReader(jsonBytes))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v", resp)
	}
}

func DownloadDataset(datasetSource string, forceDownload bool) string {

	filePath := datasetSource

	if content.IsURL(datasetSource) {
		fmt.Println("Importing dataset from URL:" + datasetSource)
		fileName := UrlToFilename(datasetSource)
		filePath := path.Join(DownloadDir, datasetDir, fileName)
		content.Download(datasetSource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided. Using local file")
	}

	return filePath
}

func MapDataset(filePath string) *model.Dataset {
	reader := content.OpenReader(filePath)

	var wdaDataset = &wda.Dataset{}
	content.ParseJson(reader, wdaDataset)

	dataset := &model.Dataset{
		ID:    wdaDataset.Ons.DatasetDetail.ID,
		Title: wdaDataset.Ons.DatasetDetail.Names.Name[0].Text,
	}

	dataset.Metadata = &model.Metadata{
		ReleaseDate: wdaDataset.Ons.DatasetDetail.PublicationDate,
		Description: mapDescription(wdaDataset.Ons.DatasetDetail.RefMetadata),
	}

	for _, wdaDimension := range wdaDataset.Ons.DatasetDetail.Dimensions.Dimension {

		// ignore hierarchies for now when mapping data sets.
		if wdaDimension.DimensionType == "Location" {
			continue
		}

		queryString := strings.Split(wdaDataset.Ons.Node.Urls.URL[0].Href, "?")[1]
		dimensionUrl := wdaDataset.Ons.Base.Href +
			path.Join("dataset",
				dataset.ID,
				"dimension",
				wdaDimension.DimensionID+".json?"+queryString)

		//fmt.Printf("wda dimension%+v\n", wdaDimension)
		dimensionPath := DownloadDimension(dimensionUrl, false)
		dimension := MapDimension(dimensionPath, wdaDimension.DimensionType)
		dataset.Dimensions = append(dataset.Dimensions, dimension)
	}

	return dataset
}

func mapDescription(refMetaData json.RawMessage) string {

	// refmetadata field can either be a single refmetadata object, or an array of them.
	// Try and unmarshall as a refmetadata object first. If it fails then attempt to unmarshall as an array.
	var metadata wda.RefMetadata
	if err := json.Unmarshal([]byte(refMetaData), &metadata); err != nil {

		var metadataArray wda.RefMetadataArray
		if err := json.Unmarshal([]byte(refMetaData), &metadataArray); err != nil {
			return string(refMetaData)
		}

		return metadataArray.RefMetadataItem[0].Descriptions.Description[0].Text
	} else {
		return metadata.RefMetadataItem.Descriptions.Description[0].Text
	}

}
