package models

type Episode struct {
	Id            int64
	Name          string
	Season        int
	EpisodeNumber *int
	Type          string
}
