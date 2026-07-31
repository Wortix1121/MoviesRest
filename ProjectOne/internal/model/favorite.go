package model

type Favorite struct {
	UserID  int64 `json:"user_id" db:"user_id"`
	MovieID int64 `json:"movie_id" db:"movie_id"`
}
