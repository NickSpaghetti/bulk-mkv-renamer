package models

type ShowSearchResult struct {
	Id       int64    `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	Language string   `json:"language,omitempty"`
	Genres   []string `json:"genres,omitempty"`
}
