package importer

import (
	"fmt"
	"path"
	"github.com/ONSdigital/dp-dd-dataset-importer/content"
	"github.com/ONSdigital/dp-dd-dataset-importer/wda"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"encoding/json"
	"log"
)

// downloadHierarchy if it does not already exist. Force a download by passing forceDownload=true
// return the filepath to the downloaded file.
func downloadHierarchy(hierarchySource string, forceDownload bool) string {
	filePath := hierarchySource

	if content.IsURL(hierarchySource) {
		fmt.Println("Importing hierarchy from URL:" + hierarchySource)
		fileName := urlToFilename(hierarchySource)
		filePath := path.Join(DownloadDir, HierarchyDir, fileName)
		content.Download(hierarchySource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided. Attempting to read file locally")
	}

	return filePath
}


func mapHierarchyToAreas(filePath string) []*model.Area {

	reader := content.OpenReader(filePath)
	var wdaHierarchy = &wda.Hierarchy{}
	content.ParseJson(reader, wdaHierarchy)

	var areas []*model.Area = mapAreas(wdaHierarchy.Ons.GeographyList.Items.Item)

	return areas
}

// mapDimensionOptionText - determine the option text format and map it as the identified type. It can either be a JSON object or an array
func mapAreas(rawArea json.RawMessage) []*model.Area {
	var wdaAreaArray wda.AreaArray
	if err := json.Unmarshal([]byte(rawArea), &wdaAreaArray); err != nil {

		var wdaArea wda.Area
		if err := json.Unmarshal([]byte(rawArea), &wdaArea); err != nil {
			log.Fatal(err)
		}

		var results []*model.Area

		area := &model.Area{
			Title:wdaArea.Labels.Label[0].Text,
			Type:wdaArea.AreaType.Codename,
		}

		results = append(results, area)
		return results

	} else {

		var results []*model.Area

		for _, wdaArea := range wdaAreaArray {

			area := &model.Area{
				Title:wdaArea.Labels.Label[0].Text,
				Type:wdaArea.AreaType.Codename,
			}
			results = append(results, area)
		}
		return results
	}
}