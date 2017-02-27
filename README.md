dp-dd-dataset-importer
================

This project was created to read data from the WDA API and convert it into a format that can be used by the dp-dd-search-indexer.

It can download dataset data from the WDA API, and convert it into the data discovery dataset JSON format. 
It can download hierarchy data and convert it to the format required by the search indexer.
It saves output locally and can also send the output to the dp-dd-search-indexer

### Getting started

#### Build

```
go build
```

#### Importing a single dataset

To import a single dataset use the -dataset flag with the WDA API URL of the dataset.
``` 
./dp-dd-dataset-importer -dataset http://web.ons.gov.uk/ons/api/data/datasetdetails/QS501EW.json?apikey={API-KEY}&context=Census&geog=2011STATH
```
The importer will save the response locally so running the application again will use the local file instead of downloading again. To override this and force the download again you can use the -force flag
``` 
./dp-dd-dataset-importer -force -dataset http://web.ons.gov.uk/ons/api/data/datasetdetails/QS501EW.json?apikey={API-KEY}&context=Census&geog=2011STATH
```
Now that the file is downloaded and available locally you can specify it directly
``` 
./dp-dd-dataset-importer -dataset downloaded/datasets/QS501EW.json
```
By default the output dataset JSON will be stored in the output directory. If you also want to put the json directly into the search indexer use the -indexer flag and pass the indexer URL
```
./dp-dd-dataset-importer -dataset downloaded/datasets/ASHE07H.json -indexer http://localhost:20050/index 
```

#### Importing multiple datasets

To import multiple datasets you can use the -datasets flag with the WDA API datasets URL
```
./dp-dd-dataset-importer -datasets http://data.ons.gov.uk/ons/api/data/datasets.json?apikey={API-KEY}
```
Once you have downloaded the datasets file once the local copy will be used. You can specify the local copy
```
./dp-dd-dataset-importer -datasets downloaded/datasets.json
```
To force another download of the datasets file you can either delete it, or run the importer with the -force flag
```
./dp-dd-dataset-importer -force -datasets http://data.ons.gov.uk/ons/api/data/datasets.json?apikey={API-KEY}
```
You can limit the number of datasets processed from each context by using the limit flag
```
./dp-dd-dataset-importer -datasets downloaded/datasets.json -limit 1
```

#### Importing hierarchy data

```
./dp-dd-dataset-importer -hierarchies http://data.ons.gov.uk/ons/api/data/hierarchies.json?apikey={API-KEY}
```
Do note that the endpoint for indexing geographic areas is different from the dataset endpoint, so adding the indexer flag will look like this:

```
-indexer http://localhost:20050/index-area
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
