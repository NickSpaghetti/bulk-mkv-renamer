package dto

type ShowSearchResultDto struct {
	Score float64 `json:"score"`
	Show  struct {
		Id             int64    `json:"id"`
		Url            string   `json:"url"`
		Name           string   `json:"name"`
		Type           string   `json:"type"`
		Language       string   `json:"language"`
		Genres         []string `json:"genres"`
		Status         string   `json:"status"`
		Runtime        int      `json:"runtime"`
		AverageRuntime int      `json:"averageRuntime"`
		Premiered      string   `json:"premiered"`
		Ended          *string  `json:"ended"`
		OfficialSite   *string  `json:"officialSite"`
		Schedule       struct {
			Time string   `json:"time"`
			Days []string `json:"days"`
		} `json:"schedule"`
		Rating struct {
			Average *float64 `json:"average"`
		} `json:"rating"`
		Weight  int `json:"weight"`
		Network struct {
			Id      int    `json:"id"`
			Name    string `json:"name"`
			Country struct {
				Name     string `json:"name"`
				Code     string `json:"code"`
				Timezone string `json:"timezone"`
			} `json:"country"`
			OfficialSite *string `json:"officialSite"`
		} `json:"network"`
		WebChannel interface{} `json:"webChannel"`
		DvdCountry interface{} `json:"dvdCountry"`
		Externals  struct {
			Tvrage  *int   `json:"tvrage"`
			Thetvdb int    `json:"thetvdb"`
			Imdb    string `json:"imdb"`
		} `json:"externals"`
		Image struct {
			Medium   string `json:"medium"`
			Original string `json:"original"`
		} `json:"image"`
		Summary string `json:"summary"`
		Updated int    `json:"updated"`
		Links   struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Previousepisode struct {
				Href string `json:"href"`
			} `json:"previousepisode"`
		} `json:"_links"`
	} `json:"show"`
}
