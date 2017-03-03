package wda

type Hierarchies struct {
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
		GeographicalHierarchyList struct {
			GeographicalHierarchy []struct {
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
				Types struct {
					GeographicalType []struct {
						XMLLang string `json:"@xml.lang"`
						Text    string `json:"$"`
					} `json:"geographicalType"`
				} `json:"types"`
				Year int `json:"year"`
			} `json:"geographicalHierarchy"`
		} `json:"geographicalHierarchyList"`
	} `json:"ons"`
}
