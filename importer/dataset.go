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
	dataset, err := mapDataset(datasetPath)
	if err != nil {
		return
	}

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
			ID:   dataset.ID,
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
func mapDataset(filePath string) (*model.Dataset, error) {
	reader := content.OpenReader(filePath)

	var wdaDataset = &wda.Dataset{}
	err := content.ParseJson(reader, wdaDataset)
	if err != nil {
		fmt.Printf("Failed to deserialise the json dataset file %v\n", filePath)
		return nil, err
	}

	dataset := &model.Dataset{
		ID:                  wdaDataset.Ons.DatasetDetail.ID,
		Title:               wdaDataset.Ons.DatasetDetail.Names.Name[0].Text,
		GeographicHierarchy: mapHierarchies(&wdaDataset.Ons.DatasetDetail.GeographicalHierarchies.GeographicalHierarchy),
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

	return dataset, nil
}

// Currently we only map a single hierarchy for a dataset. This may well need expanding to import multiple
// hierarchies for a dataset.
func mapHierarchies(wdaHierarchy *wda.HierarchySummary) []*model.GeographicHierarchySummary {
	var hierarchySummaries []*model.GeographicHierarchySummary

	hierarchy := &model.GeographicHierarchySummary{
		ID:        wdaHierarchy.ID,
		Title:     wdaHierarchy.Names.Name[0].Text,
		AreaTypes: mapAreaTypes(wdaHierarchy.AreaTypes),
	}

	hierarchySummaries = append(hierarchySummaries, hierarchy)

	return hierarchySummaries
}

// The area types section in the wda json can be either an array or an object.
// If there is more than one area type then it is represented as an array.
// If there is a single area type then it is returned as an object instead of an array with a single entry.
func mapAreaTypes(wdaAreaType json.RawMessage) []*model.AreaType {
	var mappedAreaTypes []*model.AreaType

	var areaTypesArray wda.AreaTypesArray
	if err := json.Unmarshal([]byte(wdaAreaType), &areaTypesArray); err != nil {

		var areaType wda.AreaType
		if err := json.Unmarshal([]byte(wdaAreaType), &areaType); err != nil {
			return mappedAreaTypes
		}

		mappedAreaTypes = append(mappedAreaTypes, &model.AreaType{
			Title: areaType.AreaType.Codename,
			ID:    areaType.AreaType.Abbreviation,
			Level: areaType.AreaType.Level,
		})

	} else {
		for _, areaType := range areaTypesArray.AreaType {
			mappedAreaTypes = append(mappedAreaTypes, &model.AreaType{
				Title: areaType.Codename,
				ID:    areaType.Abbreviation,
				Level: areaType.Level,
			})
		}
	}

	return mappedAreaTypes
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
