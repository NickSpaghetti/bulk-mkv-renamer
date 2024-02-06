package dto

import "time"

type EpisodeDto struct {
	Id       int64     `json:"id"`
	Url      string    `json:"url"`
	Name     string    `json:"name"`
	Season   int       `json:"season"`
	Number   *int      `json:"number,omitempty"`
	Type     string    `json:"type"`
	Airdate  string    `json:"airdate"`
	Airtime  string    `json:"airtime"`
	Airstamp time.Time `json:"airstamp"`
	Runtime  int       `json:"runtime"`
	Rating   struct {
		Average interface{} `json:"average"`
	} `json:"rating"`
	Image *struct {
		Medium   string `json:"medium"`
		Original string `json:"original"`
	} `json:"image"`
	Summary *string `json:"summary"`
	Links   struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Show struct {
			Href string `json:"href"`
		} `json:"show"`
	} `json:"_links"`
}
