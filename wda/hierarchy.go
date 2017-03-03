package wda

import "encoding/json"

type Hierarchy struct {
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
			LinkedNode struct {
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
		GeographyList struct {
			Geography struct {
				ID    string `json:"id"`
				Names struct {
					Name []struct {
						XMLLang       string `json:"@xml.lang"`
						NAMING_FAILED string `json:"$"`
					} `json:"name"`
				} `json:"names"`
			} `json:"geography"`
			Items struct {
				Item json.RawMessage `json:"item"`
			} `json:"items"`
		} `json:"geographyList"`
	} `json:"ons"`
}

type AreaArray []struct {
	Labels struct {
		Label []struct {
			XMLLang string `json:"@xml.lang"`
			Text    string `json:"$"`
		} `json:"label"`
	} `json:"labels"`
	ItemCode   string `json:"itemCode"`
	ParentCode string `json:"parentCode,omitempty"`
	AreaType   struct {
		Abbreviation string `json:"abbreviation"`
		Codename     string `json:"codename"`
		Level        int    `json:"level"`
	} `json:"areaType"`
	SubthresholdAreas string `json:"subthresholdAreas"`
}

type Area struct {
	Labels struct {
		Label []struct {
			XMLLang string `json:"@xml.lang"`
			Text    string `json:"$"`
		} `json:"label"`
	} `json:"labels"`
	ItemCode   string `json:"itemCode"`
	ParentCode string `json:"parentCode,omitempty"`
	AreaType   struct {
		Abbreviation string `json:"abbreviation"`
		Codename     string `json:"codename"`
		Level        int    `json:"level"`
	} `json:"areaType"`
	SubthresholdAreas string `json:"subthresholdAreas"`
}
