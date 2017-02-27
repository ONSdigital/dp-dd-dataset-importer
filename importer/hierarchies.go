package importer

import (
	"fmt"
	"github.com/ONSdigital/dp-dd-dataset-importer/wda"
	"github.com/ONSdigital/dp-dd-dataset-importer/content"
	"encoding/json"
	"net/http"
	"bytes"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"log"
)

func ImportHierarchies(hierarchiesSource string, forceDownload bool, indexerUrl string) {
	filePath := hierarchiesSource

	if content.IsURL(hierarchiesSource) {
		fileName := urlToFilename(hierarchiesSource)
		filePath := "./" + DownloadDir + "/" + fileName
		content.Download(hierarchiesSource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided for datasets. Attempting to read the datasets file locally")
	}

	reader := content.OpenReader(filePath)

	var hierarchies = &wda.Hierarchies{}
	content.ParseJson(reader, hierarchies)

	for _, hierarchy := range hierarchies.Ons.GeographicalHierarchyList.GeographicalHierarchy{
		fmt.Println("Importing hierarchy: " + hierarchy.Names.Name[0].Text)

		for _, url := range hierarchy.Urls.URL {

			if url.Representation == "json" {

				// the default url provided from the hierarchy list adds levels 0,1, and 2. Here we add 3 and 4.
				fullUrl := hierarchies.Ons.Base.Href + url.Href + ",3,4,5"
				//fmt.Println(fullUrl)
				filepath :=  downloadHierarchy(fullUrl, forceDownload)
				areas := mapHierarchyToAreas(filepath)

				if len(indexerUrl) > 0 {

					for _, area := range areas {
						fmt.Println("Sending document to indexer " + indexerUrl)
						document := &model.Document{
							Body: area,
							Type: "area",
						}
						jsonBytes, err := json.Marshal(document)
						if err != nil {
							log.Fatal(err)
						}
						resp, err := http.Post(indexerUrl, "application/json", bytes.NewReader(jsonBytes))
						if err != nil {
							log.Fatal(err)
						}
						fmt.Printf("%+v\n", resp)
					}


				}
			}
		}
	}
}