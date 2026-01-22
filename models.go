package main

type Product struct {
	Title     string
	Price     string
	Image     string
	Link      string
	Views     string
	User      string
	MediaType string
}

type ScrapingConfig struct {
	BaseURL        string
	PagePattern    string
	TotalPages     int
	ContainerClass string
	TitleClass     string
	ImageClass     string
	LinkClass      string
	ExtraSelectors map[string]string
}

type PageStructure struct {
	ContainerSelector string
	TitleSelector     string
	ImageSelector     string
	LinkSelector      string
	Selectors         map[string]string
}