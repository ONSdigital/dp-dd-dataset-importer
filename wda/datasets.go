package wda

type Datasets struct {
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
		DatasetList struct {
			Contexts struct {
				Context []struct {
					ContextName string `json:"contextName"`
					Datasets    struct {
						Dataset []struct {
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
							GeographicalHierarchy string `json:"geographicalHierarchy"`
							PublicationDate       string `json:"publicationDate"`
							IsHidden              string `json:"isHidden"`
							IsGeoSignificant      string `json:"isGeoSignificant"`
							Replacement           bool   `json:"replacement"`
						} `json:"dataset"`
					} `json:"datasets"`
				} `json:"context"`
			} `json:"contexts"`
		} `json:"datasetList"`
	} `json:"ons"`
}
