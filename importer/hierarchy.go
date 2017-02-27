package importer

import (
	"fmt"
	"path"
	"github.com/ONSdigital/dp-dd-dataset-importer/content"
	"github.com/ONSdigital/dp-dd-dataset-importer/wda"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
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

	var areas []*model.Area

	for _, wdaArea := range wdaHierarchy.Ons.GeographyList.Items.Item {

		area := &model.Area{
			ID:wdaArea.ItemCode,
			Title:wdaArea.Labels.Label[0].NAMING_FAILED,
			Level:wdaArea.AreaType.Level,
			Type:wdaArea.AreaType.Codename,
			Geography:wdaHierarchy.Ons.GeographyList.Geography.Names.Name[0].NAMING_FAILED,
			GeographyId:wdaHierarchy.Ons.GeographyList.Geography.ID,
		}

		fmt.Println("Area %+v", area)

		areas = append(areas, area)
	}

	//fmt.Println("areas slice length:")
	//fmt.Println(len(areas))

	return areas
}