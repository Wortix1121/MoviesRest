package model

type Season struct {
	ID           int64  `json:"id" db:"id"`
	MovieID      int64  `json:"movie_id" db:"movie_id"`
	SeasonNumber int64  `json:"season_number" db:"season_number"`
	Title        string `json:"title" db:"title"`
}
