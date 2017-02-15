package wda

import (
	"encoding/json"
	"time"
)

// The WDA dimension structure.
type Dimension struct {
	Structure struct {
		Header struct {
			ID        string    `json:"ID"`
			Test      bool      `json:"Test"`
			Truncated bool      `json:"Truncated"`
			Prepared  time.Time `json:"Prepared"`
			Sender    struct {
				ID   string `json:"@id"`
				Name struct {
					XMLLang string `json:"@xml.lang"`
					Text    string `json:"$"`
				} `json:"Name"`
				Contact struct {
					Name struct {
						XMLLang string `json:"@xml.lang"`
						Text    string `json:"$"`
					} `json:"Name"`
					Telephone string `json:"Telephone"`
				} `json:"Contact"`
			} `json:"Sender"`
			Extracted time.Time `json:"Extracted"`
		} `json:"Header"`
		CodeLists struct {
			CodeList struct {
				ID       string `json:"@id"`
				AgencyID string `json:"@agencyID"`
				Version  string `json:"@version"`
				IsFinal  string `json:"@isFinal"`
				Name     []struct {
					XMLLang string `json:"@xml.lang"`
					Text    string `json:"$"`
				} `json:"Name"`
				Code []struct {
					Value       string          `json:"@value"`
					Urn         string          `json:"@urn,omitempty"`
					ParentCode  string          `json:"@parentCode,omitempty"`
					Description json.RawMessage `json:"Description"`
					//Annotations struct {
					//	Annotation struct {
					//		AnnotationType string `json:"AnnotationType"`
					//		AnnotationText struct {
					//			XMLLang string `json:"@xml.lang"`
					//			Text string `json:"$"`
					//		} `json:"AnnotationText"`
					//	} `json:"Annotation"`
					//} `json:"Annotations,omitempty"`
				} `json:"Code"`
				//Annotations struct {
				//	Annotation struct {
				//		AnnotationType string `json:"AnnotationType"`
				//		AnnotationText struct {
				//			XMLLang string `json:"@xml.lang"`
				//			Text string `json:"$"`
				//		} `json:"AnnotationText"`
				//	} `json:"Annotation"`
				//} `json:"Annotations"`
			} `json:"CodeList"`
		} `json:"CodeLists"`
	} `json:"Structure"`
}

type DescriptionArray []struct {
	XMLLang string `json:"@xml.lang"`
	Text    string `json:"$"`
}

type Description struct {
	XMLLang string `json:"@xml.lang"`
	Text    string `json:"$"`
}
