package model

import "time"

type Episode struct {
	ID              int64     `json:"id" db:"id"`
	SeasonID        int64     `json:"season_id" db:"season_id"`
	EpisodeNumber   int64     `json:"episode_number" db:"episode_number"`
	Title           string    `json:"title" db:"title"`
	DurationMinutes time.Time `json:"duration_minutes" db:"duration_minutes"`
	VideoUrl        string    `json:"video_url" db:"video_url"`
}
