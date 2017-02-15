package importer

import (
	"fmt"
	"github.com/ONSdigital/dp-dd-dataset-importer/content"
	"github.com/ONSdigital/dp-dd-dataset-importer/wda"
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
	content.ParseJson(reader, datasets)

	var datasetIdsAlreadyProcessed []string

	for _, context := range datasets.Ons.DatasetList.Contexts.Context {
		fmt.Println("Importing datasets in WDA context: " + context.ContextName)
		count := 0

	ToNextDataset:
		for _, dataset := range context.Datasets.Dataset {

			for _, item := range datasetIdsAlreadyProcessed {
				//fmt.Println("checking if item has been processed: item " + dataset.ID + " against " + item)
				if item == dataset.ID {
					fmt.Println("already processed dataset with ID " + dataset.ID + ", ignoring this one")
					break ToNextDataset
				}
			}

			for _, url := range dataset.Urls.URL {
				if url.Representation == "json" {
					ImportDataset(datasets.Ons.Base.Href+url.Href, forceDownload, indexerUrl)
					datasetIdsAlreadyProcessed = append(datasetIdsAlreadyProcessed, dataset.ID)
				}
			}

			if limit > 0 {
				count++
				if count == limit {
					break
				}
			}
		}

	}
}
