package dto

type ListSeasonResultDto struct {
	Id           int64  `json:"id"`
	Url          string `json:"url"`
	Number       int    `json:"number"`
	Name         string `json:"name"`
	EpisodeOrder int    `json:"episodeOrder"`
	PremiereDate string `json:"premiereDate"`
	EndDate      string `json:"endDate"`
	Network      struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Country struct {
			Name     string `json:"name"`
			Code     string `json:"code"`
			Timezone string `json:"timezone"`
		} `json:"country"`
		OfficialSite interface{} `json:"officialSite"`
	} `json:"network"`
	WebChannel interface{} `json:"webChannel"`
	Image      struct {
		Medium   string `json:"medium"`
		Original string `json:"original"`
	} `json:"image"`
	Summary string `json:"summary"`
	Links   struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}
