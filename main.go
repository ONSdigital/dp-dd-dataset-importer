package main

import (
	"flag"
	"fmt"
	"github.com/ONSdigital/dp-dd-dataset-importer/importer"
)

func main() {

	fmt.Println("Starting")

	datasetSource := flag.String("dataset", "", "URL or file of a single dataset to import.")
	datasetsSource := flag.String("datasets", importer.DatasetsFile, "URL or file of datasets to import.")
	limit := flag.Int("limit", 0, "limit the number of datasets downloaded from each context")
	forceDownload := flag.Bool("force", false, "if true then always download files from WDA, else use local files")
	indexerURL := flag.String("indexer", "", "The url of the search indexer service")
	flag.Parse()

	if len(*datasetSource) > 0 {
		fmt.Println("Importing a collection of datasets from " + *datasetsSource)
		importer.ImportDataset(*datasetSource, *forceDownload, *indexerURL)
		return
	}

	fmt.Println("Importing a collection of datasets from " + *datasetsSource)
	importer.ImportDatasets(*datasetsSource, *limit, *forceDownload, *indexerURL)

	fmt.Println("Finished")

}
