package importer

import (
	"log"
	"net/url"
	"path"
	"path/filepath"
)

var datasetDir string = "datasets"
var DimensionDir string = "dimensions"
var OutputDir string = "output"
var DownloadDir string = "downloaded"
var DatasetsFile string = DownloadDir + "/datasets.json"

func UrlToFilename(sourceUrl string) string {
	url, err := url.Parse(sourceUrl)
	if err != nil {
		log.Fatal(url)
	}
	fileName := path.Base(url.EscapedPath())

	var extension = filepath.Ext(fileName)
	var name = fileName[0 : len(fileName)-len(extension)]

	// Seperate dataset files based on their geography / edition
	//geography := url.Query().Get("geog")
	//if len(geography) > 0 {
	//	name += "_" + geography
	//}

	// diff = differentiator = edition (e.g. '2015')
	//edition := url.Query().Get("diff")
	//if len(edition) > 0 {
	//	name += "_" + edition
	//}

	name += ".json"

	return name
}
