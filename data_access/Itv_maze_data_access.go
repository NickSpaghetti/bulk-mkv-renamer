package data_access

import "net/http"

type ITvMazeDataAccess interface {
	FindShowIdByName(showName string) (*http.Response, error)
	ListSeasons(showId int64) (*http.Response, error)
	ListEpisodesBySeasonId(seasonId int64) (*http.Response, error)
}
