package service

import (
	"GoMkvRn/data_access"
	"GoMkvRn/dto"
	"GoMkvRn/models"
	"encoding/json"
	"fmt"
	"os"
)

type TvMazeService struct {
	tvMazeDataAccess data_access.ITvMazeDataAccess
}

func NewTvMazeService(tvMazeDataAccess data_access.ITvMazeDataAccess) *TvMazeService {
	return &TvMazeService{
		tvMazeDataAccess: tvMazeDataAccess,
	}
}

func (t *TvMazeService) FindShowIdByName(showName string) []models.ShowSearchResult {
	if len(showName) == 0 {
		fmt.Println("Cannot search for an empty show")
		os.Exit(-1)
	}

	resp, err := t.tvMazeDataAccess.FindShowIdByName(showName)
	if err != nil {
		fmt.Printf("failed to get list of tv shows with the name %s.\n", showName)
		fmt.Println(err)
		os.Exit(-1)
	}

	var tvShowResultsDtoSlice *[]dto.ShowSearchResultDto
	err = json.NewDecoder(resp.Body).Decode(&tvShowResultsDtoSlice)
	if err != nil {
		fmt.Printf("failed to read search results for %s\n", showName)
		fmt.Println(err)
	}

	var tvShowResults []models.ShowSearchResult
	for _, showResultDto := range *tvShowResultsDtoSlice {
		tvShowResults = append(tvShowResults, models.ShowSearchResult{
			Id:       showResultDto.Show.Id,
			Name:     showResultDto.Show.Name,
			Language: showResultDto.Show.Language,
			Genres:   showResultDto.Show.Genres,
		})
	}

	return tvShowResults
}

func (t *TvMazeService) ListSeasons(showId int64) []models.Season {
	if showId <= 0 {
		fmt.Printf("Invalid show id %d\n", showId)
		os.Exit(-1)
	}

	resp, err := t.tvMazeDataAccess.ListSeasons(showId)
	if err != nil {
		fmt.Printf("failed to get list of seasons with the id %d.\n", showId)
		fmt.Println(err)
		os.Exit(-1)
	}

	var tvShowResultsDtoSlice *[]dto.ListSeasonResultDto
	err = json.NewDecoder(resp.Body).Decode(&tvShowResultsDtoSlice)
	if err != nil {
		fmt.Printf("failed to read search results for %d\n", showId)
		fmt.Println(err)
	}

	var seasons []models.Season
	for _, seasonResultDto := range *tvShowResultsDtoSlice {
		seasons = append(seasons, models.Season{
			Id:                 seasonResultDto.Id,
			Name:               seasonResultDto.Name,
			Number:             seasonResultDto.Number,
			TotalEpisodeNumber: seasonResultDto.EpisodeOrder,
		})
	}

	return seasons
}

func (t *TvMazeService) ListEpisodesBySeason(seasonId int64) []models.Episode {
	if seasonId <= 0 {
		fmt.Printf("Invalid season id %d\n", seasonId)
		os.Exit(-1)
	}

	resp, err := t.tvMazeDataAccess.ListEpisodesBySeasonId(seasonId)
	if err != nil {
		fmt.Printf("failed to get list of episodes with the season id %d.\n", seasonId)
		fmt.Println(err)
		os.Exit(-1)
	}

	var episodesDto *[]dto.EpisodeDto
	err = json.NewDecoder(resp.Body).Decode(&episodesDto)
	if err != nil {
		fmt.Printf("failed to read search results for %d\n", seasonId)
		fmt.Println(err)
	}

	var episodes []models.Episode
	for _, episodeDto := range *episodesDto {
		episodes = append(episodes, models.Episode{
			Id:            episodeDto.Id,
			Name:          episodeDto.Name,
			EpisodeNumber: episodeDto.Number,
			Season:        episodeDto.Season,
			Type:          episodeDto.Type,
		})
	}

	return episodes
}
