package model

type ParserRequest struct {
	URL string `json:"url"`
}

type ParserResponse struct {
	Version       string   `json:"versionHtml"`
	Title         string   `json:"title"`
	ListH1        []string `json:"listH1,omitempty"`
	ListH2        []string `json:"listH2,omitempty"`
	ListH3        []string `json:"listH3,omitempty"`
	ListH4        []string `json:"listH4,omitempty"`
	ListH5        []string `json:"listH5,omitempty"`
	ListH6        []string `json:"listH6,omitempty"`
	InternalLinks []*Link  `json:"internalLinks,omitempty"`
	ExternalLinks []*Link  `json:"externalLinks,omitempty"`
	Login         bool     `json:"login"`
}

type Link struct {
	Name       string `json:"name"`
	Url        string `json:"url"`
	Accessible bool   `json:"accessible"`
}
