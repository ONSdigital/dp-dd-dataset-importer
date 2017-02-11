package content

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Download(source string, destination string, forceDownload bool) {

	if _, err := os.Stat(destination); err == nil {
		if !forceDownload {
			fmt.Println("Download file " + destination + " already exists. Using local copy")
			return
		}

		os.Remove(destination)
	}

	// make parent directories
	err := os.MkdirAll(path.Dir(destination), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloading file: " + source)
	reader := OpenReader(source)
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

func SaveObjectJson(source interface{}, destination string) {

	_ = os.MkdirAll(filepath.Dir(destination), os.ModePerm)

	fmt.Println("Saving output file " + destination)
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
}

// OpenReader opens a reader on a local file or a url, as appropriate
func OpenReader(endpoint string) io.ReadCloser {
	if IsURL(endpoint) {
		response, err := http.Get(endpoint)
		if err != nil {
			fmt.Println("Error calling endpoint")
			panic(err)
		}
		return response.Body
	}

	file, err := os.Open(endpoint)
	if err != nil {
		fmt.Printf("Error opening file '%s': %s\n", endpoint, err)
		panic(err)
	}
	return file

}

// Parse parses the content from the Reader into the data object
func ParseJson(reader io.ReadCloser, data interface{}) {
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error reading body! %s\n", err)
		panic(err.Error())
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		switch t := err.(type) {
		case *json.SyntaxError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Character)"
			fmt.Printf("Invalid character at offset %v\n %s\n", t.Offset, jsn)
		case *json.UnmarshalTypeError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Type)"
			fmt.Printf("Invalid value at offset %v\n %s\n", t.Offset, jsn)
			fmt.Println("You might need to save the json file locally and validate the format - e.g. classifications with a single CodeList are known to be in an invalid format. See the readme for details.")
		default:
			fmt.Printf("Unable to unmarshal data. Error=%T, data=%s\n", t, string(body))
		}
		panic(err)
	}

}

func IsURL(file string) bool {
	return strings.HasPrefix(file, "http")
}
