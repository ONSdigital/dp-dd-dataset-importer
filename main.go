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
	"github.com/ONSdigital/dp-dd-search-indexer/model"
	"encoding/json"
	"path/filepath"
	"net/http"
	"bytes"
)

var downloadDir string = "downloaded"
var outputDir string = "output"
var datasetsFile string = downloadDir + "/datasets.json"

func main() {

	fmt.Println("Starting")

	datasetSource := flag.String("dataset", "", "URL or file of a single dataset to import.")
	datasetsSource := flag.String("datasets", datasetsFile, "URL or file of datasets to import.")
	limit := flag.Int("limit", 0, "limit the number of datasets downloaded from each context")
	forceDownload := flag.Bool("force", false, "if true then always download files from WDA, else use local files")
	indexerUrl := flag.String("indexer", "", "The url of the search indexer service")
	flag.Parse()

	if len(*datasetSource) > 0 {
		// process single dataset
		fmt.Println("Processing a single dataset: " + *datasetSource)
		dataset := GetDataset(*datasetSource, *forceDownload)

		// save
		fileName := getFilename(*datasetSource)
		outputFilePath := "./" + outputDir + "/" + fileName
		saveObjectJson(dataset, outputFilePath)

		// index
		fmt.Println(len(*indexerUrl))
		fmt.Println(*indexerUrl)
		if len(*indexerUrl) > 0 {

			fmt.Println("Sending document to indexer")
			document := &model.Document{
				Body:dataset,
				Type:"dataset",
			}
			jsonBytes, err := json.Marshal(document)
			if err != nil {
				log.Fatal(err)
			}
			resp, err := http.Post(*indexerUrl, "application/json", bytes.NewReader(jsonBytes))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v", resp)
		}

		return
	} else {
		fmt.Println("Processing a collection of datasets: " + *datasetsSource)
		ProcessDatasets(*datasetsSource, *limit, *forceDownload)
	}

	fmt.Println("Finished")

}
func GetDataset(datasetSource string, forceDownload bool) *model.Dataset {

	filePath := datasetSource

	if content.IsURL(datasetSource) {
		fileName := getFilename(datasetSource)
		filePath := "./" + downloadDir + "/" + fileName
		downloadFile(datasetSource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided. Attempting to read file locally")
	}

	reader := content.OpenReader(filePath)

	var dataset = &wda.Dataset{}
	content.Parse(reader, dataset)

	fmt.Println("Dataset ID: " + dataset.Ons.DatasetDetail.ID)

	searchDataset := &model.Dataset{
		ID:dataset.Ons.DatasetDetail.ID,
		//Description:dataset.Ons.DatasetDetail.Names.Name[0].Text,
	}

	for _, dimension := range dataset.Ons.DatasetDetail.Dimensions.Dimension {
		fmt.Println("- DimensionId: " + dimension.DimensionID)
		fmt.Println("- DimensionType: " + dimension.DimensionType)
	}

	return searchDataset

}
func ProcessDatasets(datasetsSource string, limit int, forceDownload bool) {
	filePath := datasetsSource

	if content.IsURL(datasetsSource) {
		fileName := getFilename(datasetsSource)
		filePath := "./" + downloadDir + "/" + fileName
		downloadFile(datasetsSource, filePath, forceDownload)
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
		//for _, dataset := range context.Datasets.Dataset {
		//	fmt.Println(dataset.ID + ", " +
		//		dataset.Names.Name[0].Text + ", " +
		//		dataset.GeographicalHierarchy + ", " +
		//		dataset.PublicationDate)
		//}

		count := 0

		for _, dataset := range context.Datasets.Dataset {

			for _,url := range dataset.Urls.URL {
				if url.Representation == "json" {
					fmt.Println(datasets.Ons.Base.Href + url.Href)
					fmt.Println(getFilename(url.Href))
					//GetDataset(datasets.Ons.Base.Href + url.Href)
				}
			}

			if limit > 0 {
				count ++
				if count == limit {
					break
				}
			}
		}

	}
}

func getFilename(source string) (string) {
	url, err := url.Parse(source)
	if err != nil {
		log.Fatal(url)
	}
	fileName := path.Base(url.EscapedPath())

	var extension = filepath.Ext(fileName)
	var name = fileName[0:len(fileName)-len(extension)]

	geography := url.Query().Get("geog")
	if len(geography) > 0 {
		name += "_" + geography
	}

	edition := url.Query().Get("diff")
	if len(edition) > 0 {
		name += "_" + edition
	}

	name += ".json"

	return name
}

func downloadFile(source string, destination string, forceDownload bool) {

	if _, err := os.Stat(destination); err == nil {
		if !forceDownload {
			fmt.Println("File " + destination + " already exists. Using local copy")
			return
		}

		os.Remove(destination)
	}

	fmt.Println("Downloading file")
	reader := content.OpenReader(source)
	defer reader.Close()

	file, err := os.Create(destination)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Download complete")
}


func saveObjectJson(source interface{}, destination string) {

	_ = os.MkdirAll(filepath.Dir(destination), os.ModePerm)

	fmt.Println("Saving file " + destination)
	if _, err := os.Stat(destination); err == nil {
		os.Remove(destination)
	}
	file, err := os.Create(destination)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(source)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Save complete")

}
