package model

import "time"

type Movie struct {
	ID              int64     `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	DurationMinutes int64     `json:"duration_minutes" db:"duration_minutes"`
	ReleaseDate     time.Time `json:"release_date" db:"release_date"`
	PosterUrl       string    `json:"poster_url" db:"poster_url"`
	AgeRating       string    `json:"age_rating" db:"age_rating"`
	Type            string    `json:"type" db:"type"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
