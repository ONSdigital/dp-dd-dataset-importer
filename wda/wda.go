package wda

import (
	"time"
	"encoding/json"
)

type Datasets struct {
	Ons struct {
		Base struct {
			Href string `json:"@href"`
		} `json:"base"`
		Node struct {
			Urls struct {
				URL []struct {
					Representation string `json:"@representation"`
					Href string `json:"href"`
				} `json:"url"`
			} `json:"urls"`
			Description string `json:"description"`
			Name string `json:"name"`
		} `json:"node"`
		LinkedNodes struct {
			LinkedNode struct {
				Urls struct {
					URL []struct {
						Representation string `json:"@representation"`
						Href string `json:"href"`
					} `json:"url"`
				} `json:"urls"`
				Name string `json:"name"`
				Relation string `json:"relation"`
			} `json:"linkedNode"`
		} `json:"linkedNodes"`
		DatasetList struct {
			Contexts struct {
				Context []struct {
					ContextName string `json:"contextName"`
					Datasets struct {
						Dataset []struct {
							ID string `json:"id"`
							Names struct {
								Name []struct {
									XMLLang string `json:"@xml.lang"`
									Text string `json:"$"`
								} `json:"name"`
							} `json:"names"`
							Urls struct {
								URL []struct {
									Representation string `json:"@representation"`
									Href string `json:"href"`
								} `json:"url"`
							} `json:"urls"`
							GeographicalHierarchy string `json:"geographicalHierarchy"`
							PublicationDate string `json:"publicationDate"`
							IsHidden string `json:"isHidden"`
							IsGeoSignificant string `json:"isGeoSignificant"`
							Replacement bool `json:"replacement"`
						} `json:"dataset"`
					} `json:"datasets"`
				} `json:"context"`
			} `json:"contexts"`
		} `json:"datasetList"`
	} `json:"ons"`
}

type Dataset struct {
	Ons struct {
		Base struct {
			Href string `json:"@href"`
		} `json:"base"`
		Node struct {
			Urls struct {
				URL []struct {
					Representation string `json:"@representation"`
					Href string `json:"href"`
				} `json:"url"`
			} `json:"urls"`
			Description string `json:"description"`
			Name string `json:"name"`
		} `json:"node"`
		LinkedNodes struct {
			LinkedNode []struct {
				Urls struct {
					URL []struct {
						Representation string `json:"@representation"`
						Href string `json:"href"`
					} `json:"url"`
				} `json:"urls"`
				Name string `json:"name"`
				Relation string `json:"relation"`
			} `json:"linkedNode"`
		} `json:"linkedNodes"`
		DatasetDetail struct {
			ID string `json:"id"`
			Names struct {
				Name []struct {
					XMLLang string `json:"@xml.lang"`
					Text string `json:"$"`
				} `json:"name"`
			} `json:"names"`
			Urls struct {
				URL []struct {
					Representation string `json:"@representation"`
					Href string `json:"href"`
				} `json:"url"`
			} `json:"urls"`
			RefMetadata json.RawMessage `json:"refMetadata"`
			Dimensions struct {
				Dimension []struct {
					DimensionID string `json:"dimensionId"`
					DimensionTitles struct {
						DimensionTitle []struct {
							XMLLang string `json:"@xml.lang"`
							Text string `json:"$"`
						} `json:"dimensionTitle"`
					} `json:"dimensionTitles"`
					DimensionType string `json:"dimensionType"`
					NumberOfDimensionItems int `json:"numberOfDimensionItems"`
					AdvisoryNote string `json:"advisoryNote"`
				} `json:"dimension"`
			} `json:"dimensions"`
			Designations struct {
				Designation []struct {
					XMLLang string `json:"@xml.lang"`
					Text string `json:"$"`
				} `json:"designation"`
			} `json:"designations"`
			Geocoverages struct {
				Geocoverage []struct {
					XMLLang string `json:"@xml.lang"`
					Text string `json:"$"`
				} `json:"geocoverage"`
			} `json:"geocoverages"`
			PublicationDate string `json:"publicationDate"`
			Contact struct {
				SystemID string `json:"systemId"`
				ContactName string `json:"contactName"`
				ContactEmail string `json:"contactEmail"`
				ContactPhoneNumber string `json:"contactPhoneNumber"`
			} `json:"contact"`
			StatisticalPopulations string `json:"statisticalPopulations"`
			GeographicalHierarchies struct {
				GeographicalHierarchy struct {
					ID string `json:"id"`
					Names struct {
						Name []struct {
							XMLLang string `json:"@xml.lang"`
							Text string `json:"$"`
						} `json:"name"`
					} `json:"names"`
					Urls struct {
						URL []struct {
							Representation string `json:"@representation"`
							Href string `json:"href"`
						} `json:"url"`
					} `json:"urls"`
					Time int `json:"time"`
					Types struct {
						GeographicalType []struct {
							XMLLang string `json:"@xml.lang"`
							Text string `json:"$"`
						} `json:"geographicalType"`
					} `json:"types"`
					//Differentiators struct {
					//	Differentiator string `json:"differentiator"`
					//} `json:"differentiators"`
					AreaTypes struct {
						AreaType []struct {
							Abbreviation string `json:"abbreviation"`
							Codename string `json:"codename"`
							Level int `json:"level"`
						} `json:"areaType"`
					} `json:"areaTypes"`
				} `json:"geographicalHierarchy"`
			} `json:"geographicalHierarchies"`
			IsHidden string `json:"isHidden"`
			IsGeoSignificant string `json:"isGeoSignificant"`
			IsSparse string `json:"isSparse"`
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
			ObsCount int `json:"obsCount"`
			SuppressMap string `json:"suppressMap"`
			SuppressChart string `json:"suppressChart"`
			SuppressView string `json:"suppressView"`
		} `json:"datasetDetail"`
	} `json:"ons"`
}

type RefMetadata struct {
	RefMetadataItem struct {
		Type string `json:"type"`
		SystemID string `json:"systemId"`
		Descriptions struct {
			Description []struct {
				XMLLang string `json:"@xml.lang"`
				Text string `json:"$"`
			} `json:"description"`
		} `json:"descriptions"`
		DisplayOrder int `json:"displayOrder"`
	} `json:"refMetadataItem"`
}

type RefMetadataArray struct {
	RefMetadataItem []struct {
		Type string `json:"type"`
		SystemID string `json:"systemId"`
		Descriptions struct {
			Description []struct {
				XMLLang string `json:"@xml.lang"`
				Text string `json:"$"`
			} `json:"description"`
		} `json:"descriptions"`
		DisplayOrder int `json:"displayOrder"`
	} `json:"refMetadataItem"`
}

type Dimension struct {
	Structure struct {
		Header struct {
			ID string `json:"ID"`
			Test bool `json:"Test"`
			Truncated bool `json:"Truncated"`
			Prepared time.Time `json:"Prepared"`
			Sender struct {
				ID string `json:"@id"`
				Name struct {
					XMLLang string `json:"@xml.lang"`
					Text string `json:"$"`
				} `json:"Name"`
				Contact struct {
					Name struct {
						XMLLang string `json:"@xml.lang"`
						Text string `json:"$"`
					} `json:"Name"`
					Telephone string `json:"Telephone"`
				} `json:"Contact"`
			} `json:"Sender"`
			Extracted time.Time `json:"Extracted"`
		} `json:"Header"`
		CodeLists struct {
			CodeList struct {
				ID string `json:"@id"`
				AgencyID string `json:"@agencyID"`
				Version string `json:"@version"`
				IsFinal string `json:"@isFinal"`
				Name []struct {
					XMLLang string `json:"@xml.lang"`
					Text string `json:"$"`
				} `json:"Name"`
				Code []struct {
					Value string `json:"@value"`
					Urn string `json:"@urn,omitempty"`
					ParentCode string `json:"@parentCode,omitempty"`
					Description []struct {
						XMLLang string `json:"@xml.lang"`
						Text string `json:"$"`
					} `json:"Description"`
					Annotations struct {
						Annotation struct {
							AnnotationType string `json:"AnnotationType"`
							AnnotationText struct {
								XMLLang string `json:"@xml.lang"`
								Text string `json:"$"`
							} `json:"AnnotationText"`
						} `json:"Annotation"`
					} `json:"Annotations,omitempty"`
				} `json:"Code"`
				Annotations struct {
					Annotation struct {
						AnnotationType string `json:"AnnotationType"`
						AnnotationText struct {
							XMLLang string `json:"@xml.lang"`
							Text string `json:"$"`
						} `json:"AnnotationText"`
					} `json:"Annotation"`
				} `json:"Annotations"`
			} `json:"CodeList"`
		} `json:"CodeLists"`
	} `json:"Structure"`
}