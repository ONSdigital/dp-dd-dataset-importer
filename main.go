package main

import (
	"os"
	"log"
	"fmt"
	"flag"
	"github.com/ONSdigital/dp-dd-dataset-importer/content"
	"io"
	"net/url"
	"path"
	"github.com/ONSdigital/dp-dd-dataset-importer/wda"
)

var downloadDir string = "downloaded"
var datasetsFile string = downloadDir + "/datasets.json"

func main() {

	fmt.Println("Starting")

	datasetSource := flag.String("dataset", "", "URL or file of a single dataset to import.")
	datasetsSource := flag.String("datasets", datasetsFile, "URL or file of datasets to import.")
	flag.Parse()

	if len(*datasetSource) > 0 {
		// process single dataset
		fmt.Println("Processing a single dataset: " + *datasetSource)
		ProcessDataset(*datasetSource)
		return
	} else {
		fmt.Println("Processing a collection of datasets: " + *datasetsSource)
		ProcessDatasets(*datasetsSource)
	}

	fmt.Println("Finished")

}
func ProcessDataset(datasetSource string) {

	filePath := datasetSource

	if content.IsURL(datasetSource) {
		filePath = downloadFileIfNotExists(datasetSource)
	} else {
		fmt.Println("URL was not provided. Attempting to read file locally")
	}

	reader := content.OpenReader(filePath)

	var dataset = &wda.Dataset{}
	content.Parse(reader, dataset)

	fmt.Println("Dataset ID: " + dataset.Ons.DatasetDetail.ID)
	for _, dimension := range dataset.Ons.DatasetDetail.Dimensions.Dimension {
		fmt.Println("- DimensionId: " + dimension.DimensionID)
		fmt.Println("- DimensionType: " + dimension.DimensionType)
	}
}
func ProcessDatasets(datasetsSource string) {
	filePath := datasetsSource

	if content.IsURL(datasetsSource) {
		filePath = downloadFileIfNotExists(datasetsSource)
	} else {
		fmt.Println("URL was not provided. Attempting to read file locally")
	}

	reader := content.OpenReader(filePath)

	var datasets = &wda.Datasets{}
	content.Parse(reader, datasets)

	for _, context := range datasets.Ons.DatasetList.Contexts.Context {
		fmt.Println(context.ContextName)

		//fmt.Printf("Datasets: %v", len(context.Datasets.Dataset))
		//fmt.Println()
		for _, dataset := range context.Datasets.Dataset {
			fmt.Println(dataset.ID + ", " +
				dataset.Names.Name[0].Text + ", " +
				dataset.GeographicalHierarchy + ", " +
				dataset.PublicationDate)
		}
	}
}

// download the file and store it with the filename in the url.
func downloadFileIfNotExists(url string) (string) {
	fileName := getFilename(url)
	filePath := "./downloaded/" + fileName

	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("File " + filePath + " already exists. Using local copy")
	} else {
		downloadFileToDestination(url, filePath)
	}
	return filePath
}

func getFilename(source string) (string) {
	url, err := url.Parse(source)
	if err != nil {
		log.Fatal(url)
	}
	fileName := path.Base(url.EscapedPath())
	return fileName
}

func downloadFileToDestination(source string, destination string) {

	fmt.Println("Downloading file")
	reader := content.OpenReader(source)
	defer reader.Close()
	if _, err := os.Stat(destination); err == nil {
		os.Remove(destination)
	}
	file, err := os.Create(destination)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Download complete")

}
