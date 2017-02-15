package importer

import (
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-dd-dataset-importer/content"
	"github.com/ONSdigital/dp-dd-dataset-importer/wda"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"log"
	"path"
)

// downloadDimension if it does not already exist. Force a download by passing forceDownload=true
// return the filepath to the downloaded file.
func downloadDimension(dimensionSource string, forceDownload bool) string {
	filePath := dimensionSource

	if content.IsURL(dimensionSource) {
		fmt.Println("Importing dimension from URL:" + dimensionSource)
		fileName := urlToFilename(dimensionSource)
		filePath := path.Join(DownloadDir, DimensionDir, fileName)
		content.Download(dimensionSource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided. Attempting to read file locally")
	}

	return filePath
}

// mapDimension from the WDA API format to the model.Dimension structure.
// the dimension type is only provided at the dataset level, hence it being passed in instead of mapped.
func mapDimension(filePath string, dimensionType string) *model.Dimension {

	reader := content.OpenReader(filePath)

	var wdaDimension = &wda.Dimension{}
	content.ParseJson(reader, wdaDimension)

	//fmt.Println("Dimension ID: " + wdaDimension.Structure.Header.ID)
	dimension := &model.Dimension{
		ID:   wdaDimension.Structure.CodeLists.CodeList.ID,
		Name: wdaDimension.Structure.CodeLists.CodeList.Name[0].Text,
		Type: dimensionType,
	}

	for _, wdaDimensionOption := range wdaDimension.Structure.CodeLists.CodeList.Code {

		// all dimensions in WDA seem to have a "name": "Not Applicable" option. Ignoring these.
		if wdaDimensionOption.Value == "_Z" {
			continue
		}

		var optionName string = mapDimensionOptionText(wdaDimensionOption.Description)

		dimension.Options = append(dimension.Options, &model.DimensionOption{
			ID:   wdaDimensionOption.Value,
			Name: optionName,
		})
	}

	return dimension
}

// mapDimensionOptionText - determine the option text format and map it as the identified type. It can either be a JSON object or an array
func mapDimensionOptionText(optionDescription json.RawMessage) string {
	var description wda.Description
	if err := json.Unmarshal([]byte(optionDescription), &description); err != nil {
		var descriptionArray wda.DescriptionArray
		if err := json.Unmarshal([]byte(optionDescription), &descriptionArray); err != nil {
			log.Fatal(err)
		}

		return descriptionArray[0].Text
	} else {
		return description.Text
	}
}
