package importer

import (
	"fmt"
	"github.com/ONSdigital/dp-dd-dataset-importer/content"
	"github.com/ONSdigital/dp-dd-dataset-importer/wda"
	"log"
)

// ImportDatasets - Read a list of datasets from WDA using the given datasetsSource and save to local disk.
// A smaller number of datasets can be imported by specifying the limit > 0
// If the file has already been downloaded, the local copy will be used instead unless forceDownload is true.
// Map it to the DD dataset json model and save to local disk.
// If the indexerUrl is provided then the result will be sent to it.
func ImportDatasets(datasetsSource string, limit int, forceDownload bool, indexerUrl string) {
	filePath := datasetsSource

	if content.IsURL(datasetsSource) {
		fileName := urlToFilename(datasetsSource)
		filePath := "./" + DownloadDir + "/" + fileName
		content.Download(datasetsSource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided for datasets. Attempting to read the datasets file locally")
	}

	reader := content.OpenReader(filePath)

	var datasets = &wda.Datasets{}
	err := content.ParseJson(reader, datasets)
	if err != nil {
		fmt.Println("Failed to deserialise json for the list of datasets.")
		return
	}

	var datasetIdsAlreadyProcessed map[string]struct{} = make(map[string]struct{})

	for _, context := range datasets.Ons.DatasetList.Contexts.Context {
		fmt.Println("Importing datasets in WDA context: " + context.ContextName)
		count := 0

		log.Printf("Number of datasets %v \n", len(context.Datasets.Dataset))

		for _, dataset := range context.Datasets.Dataset {

			log.Printf("dataset id: %v", dataset.ID)

			if _, exists := datasetIdsAlreadyProcessed[dataset.ID]; exists {
				fmt.Println("already processed dataset with ID " + dataset.ID + ", ignoring this one")
				continue
			}

			for _, url := range dataset.Urls.URL {
				if url.Representation == "json" {
					fmt.Println("Processing dataset " + dataset.ID + "")
					ImportDataset(datasets.Ons.Base.Href+url.Href, forceDownload, indexerUrl)
					datasetIdsAlreadyProcessed[dataset.ID] = struct{}{}
					log.Printf("datasets processed: %v", datasetIdsAlreadyProcessed)

				}
			}

			if limit > 0 {
				count++
				if count == limit {
					fmt.Println("hit the limit ")
					break
				}
			}
		}

	}
}
