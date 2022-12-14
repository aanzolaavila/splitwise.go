package resources

type CategoryID Identifier

type Category struct {
	ID        CategoryID `json:"id"`
	Name      string     `json:"name"`
	Icon      string     `json:"icon"`
	IconTypes struct {
		Slim struct {
			Small string `json:"small"`
			Large string `json:"large"`
		} `json:"slim"`
		Square struct {
			Large  string `json:"large"`
			Xlarge string `json:"xlarge"`
		} `json:"square"`
	} `json:"icon_types"`
}

type MainCategory struct {
	Category
	Subcategories []Category `json:"subcategories"`
}
