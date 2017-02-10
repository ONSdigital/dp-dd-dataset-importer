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
var datasetDir string = "datasets"
var dimensionDir string = "dimensions"

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
		ProcessDataset(*datasetSource, *forceDownload, *indexerUrl)
		return
	} else {
		fmt.Println("Processing a collection of datasets: " + *datasetsSource)
		ProcessDatasets(*datasetsSource, *limit, *forceDownload, *indexerUrl)
	}

	fmt.Println("Finished")

}

func ProcessDataset(datasetSource string, forceDownload bool, indexerUrl string) {

	fmt.Println("Processing a single dataset: " + datasetSource)
	datasetPath := DownloadDataset(datasetSource, forceDownload)
	dataset := MapDataset(datasetPath)

	// save
	fileName := getFilename(datasetSource)
	outputFilePath :=  path.Join(outputDir, fileName)
	saveObjectJson(dataset, outputFilePath)

	// index
	fmt.Println(len(indexerUrl))
	fmt.Println(indexerUrl)
	if len(indexerUrl) > 0 {

		fmt.Println("Sending document to indexer")
		document := &model.Document{
			Body:dataset,
			Type:"dataset",
		}
		jsonBytes, err := json.Marshal(document)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post(indexerUrl, "application/json", bytes.NewReader(jsonBytes))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v", resp)
	}
}

func DownloadDataset(datasetSource string, forceDownload bool) string {

	filePath := datasetSource

	if content.IsURL(datasetSource) {
		fileName := getFilename(datasetSource)
		filePath := path.Join(downloadDir, datasetDir, fileName)
		fmt.Println("download filepath:" + filePath)
		downloadFile(datasetSource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided. Using local file")
	}

	return filePath
}

func MapDataset(filePath string) *model.Dataset {
	reader := content.OpenReader(filePath)

	var wdaDataset = &wda.Dataset{}
	content.Parse(reader, wdaDataset)

	fmt.Println("Dataset ID: " + wdaDataset.Ons.DatasetDetail.ID)

	dataset := &model.Dataset{
		ID:wdaDataset.Ons.DatasetDetail.ID,
		//Description:wdaDataset.Ons.DatasetDetail.Names.Name[0].Text,
	}

	// refmetadata field can either be a single refmetadata object, or an array of them.
	// Try and unmarshall as a refmetadata object first. If it fails then attempt to unmarshall as an array.
	var metadata wda.RefMetadata
	if err := json.Unmarshal([]byte(wdaDataset.Ons.DatasetDetail.RefMetadata), &metadata); err != nil {
		var metadataArray wda.RefMetadataArray
		if err := json.Unmarshal([]byte(wdaDataset.Ons.DatasetDetail.RefMetadata), &metadataArray); err != nil {
			log.Fatal(err)
		}

		dataset.Description = metadataArray.RefMetadataItem[0].Descriptions.Description[0].Text
		fmt.Println("found a metadata slice: " + dataset.Description)
	} else {
		dataset.Description = metadata.RefMetadataItem.Descriptions.Description[0].Text
		fmt.Println("found a metadata field " + dataset.Description)
	}

	for _, dimension := range wdaDataset.Ons.DatasetDetail.Dimensions.Dimension {
		fmt.Println("- DimensionId: " + dimension.DimensionID)
		fmt.Println("- DimensionType: " + dimension.DimensionType)
	}

	return dataset
}

func DownloadDimension(dimensionSource string, forceDownload bool) string {
	filePath := dimensionSource

	if content.IsURL(dimensionSource) {
		fileName := getFilename(dimensionSource)
		filePath := path.Join(downloadDir, dimensionDir, fileName)
		fmt.Printf("download filepath:" + filePath)
		downloadFile(dimensionSource, filePath, forceDownload)
	} else {
		fmt.Println("URL was not provided. Attempting to read file locally")
	}

	return filePath
}

func MapDimension(filePath string) *model.Dimension {

	reader := content.OpenReader(filePath)

	var wdaDimension = &wda.Dimension{}
	content.Parse(reader, wdaDimension)

	fmt.Println("Dimension ID: " + wdaDimension.Structure.Header.ID)

	dimension := &model.Dimension{
		ID:wdaDimension.Structure.Header.ID,
		//Description:dataset.Ons.DatasetDetail.Names.Name[0].Text,
	}

	//for _, dimension := range dataset.Ons.DatasetDetail.Dimensions.Dimension {
	//	fmt.Println("- DimensionId: " + dimension.DimensionID)
	//	fmt.Println("- DimensionType: " + dimension.DimensionType)
	//}

	return dimension
}


func ProcessDatasets(datasetsSource string, limit int, forceDownload bool, indexerUrl string) {
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
					//fmt.Println(datasets.Ons.Base.Href + url.Href)
					//fmt.Println(getFilename(url.Href))
					ProcessDataset(datasets.Ons.Base.Href + url.Href, forceDownload, indexerUrl)
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

	// make parent directories
	err := os.MkdirAll(path.Dir(destination), os.ModePerm)
	if err != nil {
		log.Fatal(err)
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
