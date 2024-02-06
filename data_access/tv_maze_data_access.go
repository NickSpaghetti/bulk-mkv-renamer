package data_access

import (
	"fmt"
	"net/http"
	"net/url"
)

type TvMazeDataAccess struct {
	ITvMazeDataAccess
}

func (t TvMazeDataAccess) FindShowIdByName(showName string) (*http.Response, error) {
	params := url.Values{}
	params.Add("q", showName)
	uri := fmt.Sprintf("https://api.tvmaze.com/search/shows?%s", params.Encode())
	resp, err := http.Get(uri)
	return resp, err
}

func (t TvMazeDataAccess) ListSeasons(showId int64) (*http.Response, error) {
	uri := fmt.Sprintf("https://api.tvmaze.com/shows/%d/seasons", showId)
	resp, err := http.Get(uri)
	return resp, err
}
func (t TvMazeDataAccess) ListEpisodesBySeasonId(seasonId int64) (*http.Response, error) {
	uri := fmt.Sprintf("https://api.tvmaze.com/seasons/%d/episodes", seasonId)
	resp, err := http.Get(uri)
	return resp, err
}
