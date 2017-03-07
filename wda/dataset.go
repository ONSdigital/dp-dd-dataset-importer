package wda

import "encoding/json"

// The WDA dataset structure.
type Dataset struct {
	Ons struct {
		Base struct {
			Href string `json:"@href"`
		} `json:"base"`
		Node struct {
			Urls struct {
				URL []struct {
					Representation string `json:"@representation"`
					Href           string `json:"href"`
				} `json:"url"`
			} `json:"urls"`
			Description string `json:"description"`
			Name        string `json:"name"`
		} `json:"node"`
		LinkedNodes struct {
			LinkedNode []struct {
				Urls struct {
					URL []struct {
						Representation string `json:"@representation"`
						Href           string `json:"href"`
					} `json:"url"`
				} `json:"urls"`
				Name     string `json:"name"`
				Relation string `json:"relation"`
			} `json:"linkedNode"`
		} `json:"linkedNodes"`
		DatasetDetail struct {
			ID    string `json:"id"`
			Names struct {
				Name []struct {
					XMLLang string `json:"@xml.lang"`
					Text    string `json:"$"`
				} `json:"name"`
			} `json:"names"`
			Urls struct {
				URL []struct {
					Representation string `json:"@representation"`
					Href           string `json:"href"`
				} `json:"url"`
			} `json:"urls"`
			RefMetadata json.RawMessage `json:"refMetadata"`
			Dimensions  struct {
				Dimension []struct {
					DimensionID     string `json:"dimensionId"`
					DimensionTitles struct {
						DimensionTitle []struct {
							XMLLang string `json:"@xml.lang"`
							Text    string `json:"$"`
						} `json:"dimensionTitle"`
					} `json:"dimensionTitles"`
					DimensionType          string `json:"dimensionType"`
					NumberOfDimensionItems int    `json:"numberOfDimensionItems"`
					AdvisoryNote           string `json:"advisoryNote"`
				} `json:"dimension"`
			} `json:"dimensions"`
			Designations struct {
				Designation []struct {
					XMLLang string `json:"@xml.lang"`
					Text    string `json:"$"`
				} `json:"designation"`
			} `json:"designations"`
			Geocoverages struct {
				Geocoverage []struct {
					XMLLang string `json:"@xml.lang"`
					Text    string `json:"$"`
				} `json:"geocoverage"`
			} `json:"geocoverages"`
			PublicationDate string `json:"publicationDate"`
			Contact         struct {
				SystemID           string `json:"systemId"`
				ContactName        string `json:"contactName"`
				ContactEmail       string `json:"contactEmail"`
				ContactPhoneNumber string `json:"contactPhoneNumber"`
			} `json:"contact"`
			StatisticalPopulations  string `json:"statisticalPopulations"`
			GeographicalHierarchies struct {
				GeographicalHierarchy HierarchySummary `json:"geographicalHierarchy"`
			} `json:"geographicalHierarchies"`
			IsHidden         string `json:"isHidden"`
			IsGeoSignificant string `json:"isGeoSignificant"`
			IsSparse         string `json:"isSparse"`
			//Documents struct {
			//	Document []struct {
			//		Type string `json:"@type"`
			//		Href struct {
			//			XMLLang string `json:"@xml.lang"`
			//			Text string `json:"$"`
			//		} `json:"href"`
			//		Filesize struct {
			//			XMLLang string `json:"@xml.lang"`
			//			Text string `json:"$"`
			//		} `json:"filesize"`
			//	} `json:"document"`
			//} `json:"documents"`
			ObsCount      int    `json:"obsCount"`
			SuppressMap   string `json:"suppressMap"`
			SuppressChart string `json:"suppressChart"`
			SuppressView  string `json:"suppressView"`
		} `json:"datasetDetail"`
	} `json:"ons"`
}

// A dataset contains only summary information about the hierarchy.
type HierarchySummary struct {
	ID    string `json:"id"`
	Names struct {
		Name []struct {
			XMLLang string `json:"@xml.lang"`
			Text    string `json:"$"`
		} `json:"name"`
	} `json:"names"`
	Urls struct {
		URL []struct {
			Representation string `json:"@representation"`
			Href           string `json:"href"`
		} `json:"url"`
	} `json:"urls"`
	Time  int `json:"time"`
	Types struct {
		GeographicalType []struct {
			XMLLang string `json:"@xml.lang"`
			Text    string `json:"$"`
		} `json:"geographicalType"`
	} `json:"types"`
	//Differentiators struct {
	//	Differentiator string `json:"differentiator"`
	//} `json:"differentiators"`
	AreaTypes json.RawMessage `json:"areaTypes"`
}

// The area type is typically an array, but if there is only a single area type
// then it is provided as an object instead of an array with a single item.
type AreaTypesArray struct {
	AreaType []struct {
		Abbreviation string `json:"abbreviation"`
		Codename     string `json:"codename"`
		Level        int    `json:"level"`
	} `json:"areaType"`
}

type AreaType struct {
	AreaType struct {
		Abbreviation string `json:"abbreviation"`
		Codename     string `json:"codename"`
		Level        int    `json:"level"`
	} `json:"areaType"`
}

type RefMetadata struct {
	RefMetadataItem struct {
		Type         string `json:"type"`
		SystemID     string `json:"systemId"`
		Descriptions struct {
			Description []struct {
				XMLLang string `json:"@xml.lang"`
				Text    string `json:"$"`
			} `json:"description"`
		} `json:"descriptions"`
		DisplayOrder int `json:"displayOrder"`
	} `json:"refMetadataItem"`
}

type RefMetadataArray struct {
	RefMetadataItem []struct {
		Type         string `json:"type"`
		SystemID     string `json:"systemId"`
		Descriptions struct {
			Description []struct {
				XMLLang string `json:"@xml.lang"`
				Text    string `json:"$"`
			} `json:"description"`
		} `json:"descriptions"`
		DisplayOrder int `json:"displayOrder"`
	} `json:"refMetadataItem"`
}
