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

// Read a dataset from WDA using the given datasetSource and save to local disk.
// If the file has already been downloaded, the local copy will be used instead unless forceDownload is true.
// Map it to the DD dataset json model and save to local disk.
// If the indexerUrl is provided then the result will be sent to it.
func ImportDataset(datasetSource string, forceDownload bool, indexerUrl string) {

	fmt.Println("Importing dataset: " + datasetSource)
	datasetPath := downloadDataset(datasetSource, forceDownload)
	dataset := mapDataset(datasetPath)

	// save
	fileName := urlToFilename(datasetSource)
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

// downloadDataset if it does not already exist. Force a download by passing forceDownload=true
// return the filepath to the downloaded file.
func downloadDataset(datasetSource string, forceDownload bool) string {

	filePath := datasetSource

	if content.IsURL(datasetSource) {
		fmt.Println("Importing dataset from URL:" + datasetSource)
		fileName := urlToFilename(datasetSource)
		filePath := path.Join(DownloadDir, datasetDir, fileName)
		content.Download(datasetSource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided. Using local file")
	}

	return filePath
}

// mapDataset from the given filePath in the WDA format to the model.Dataset structure.
func mapDataset(filePath string) *model.Dataset {
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
		dimensionPath := downloadDimension(dimensionUrl, false)
		dimension := mapDimension(dimensionPath, wdaDimension.DimensionType)
		dataset.Dimensions = append(dataset.Dimensions, dimension)
	}

	return dataset
}

// mapDescription handles the dynamic format of the WDA refMetaData field. Sometimes its an array and others its a single object.
// If it cannot map the field as an object it attempts to map it as an array.
func mapDescription(refMetaData json.RawMessage) string {

	// refmetadata field can either be a single refmetadata object, or an array of them.
	// Try and unmarshall as a refmetadata object first. If it fails then attempt to unmarshall as an array.
	var metadata wda.RefMetadata
	if err := json.Unmarshal([]byte(refMetaData), &metadata); err != nil {

		var metadataArray wda.RefMetadataArray
		if err := json.Unmarshal([]byte(refMetaData), &metadataArray); err != nil {
			return string(refMetaData) // if all else fails just return it as a string.
		}

		return metadataArray.RefMetadataItem[0].Descriptions.Description[0].Text
	} else {
		return metadata.RefMetadataItem.Descriptions.Description[0].Text
	}
}
