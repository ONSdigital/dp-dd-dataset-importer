package importer

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-dd-dataset-importer/content"
	"github.com/ONSdigital/dp-dd-dataset-importer/wda"
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"log"
	"net/http"
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
	err := content.ParseJson(reader, hierarchies)
	if err != nil {
		fmt.Println("Failed to deserialise the json for the hierarchy list.")
	}

	for _, hierarchy := range hierarchies.Ons.GeographicalHierarchyList.GeographicalHierarchy {
		fmt.Println("Importing hierarchy: " + hierarchy.Names.Name[0].Text)

		for _, url := range hierarchy.Urls.URL {

			if url.Representation == "json" {

				// the default url provided from the hierarchy list adds levels 0,1, and 2. Here we add 3 and 4.
				fullUrl := hierarchies.Ons.Base.Href + url.Href + ",3"
				//fmt.Println(fullUrl)
				filepath := downloadHierarchy(fullUrl, forceDownload)
				areas := mapHierarchyToAreas(filepath)

				if len(indexerUrl) > 0 {

					for _, area := range areas {

						id := createHash(area.Title + area.Type)

						fmt.Println("Sending document to indexer " + area.Title + " " + area.Type + " - hash: " + id)
						document := &model.Document{
							ID:   id,
							Body: area,
							Type: "area",
						}
						jsonBytes, err := json.Marshal(document)
						if err != nil {
							log.Fatal(err)
						}
						_, err = http.Post(indexerUrl, "application/json", bytes.NewReader(jsonBytes))
						if err != nil {
							log.Fatal(err)
						}
						//fmt.Printf("%+v\n", resp)
					}

				}
			}
		}
	}
}

func createHash(input string) string {

	h := sha1.New()
	h.Write([]byte(input))

	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(input)
	fmt.Printf("%x\n", sha)

	return sha
}
