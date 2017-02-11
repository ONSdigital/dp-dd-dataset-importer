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

func DownloadDimension(dimensionSource string, forceDownload bool) string {
	filePath := dimensionSource

	if content.IsURL(dimensionSource) {
		fmt.Println("Importing dimension from URL:" + dimensionSource)
		fileName := UrlToFilename(dimensionSource)
		filePath := path.Join(DownloadDir, DimensionDir, fileName)
		content.Download(dimensionSource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided. Attempting to read file locally")
	}

	return filePath
}

func MapDimension(filePath string, dimensionType string) *model.Dimension {

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
