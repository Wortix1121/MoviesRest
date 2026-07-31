package model

import "time"

type WatchHistory struct {
	ID              int64     `json:"id" db:"id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	ContentType     string    `json:"content_type" db:"content_type"`
	ContentID       int64     `json:"content_id" db:"content_id"`
	ProgressSeconds int64     `json:"progress_seconds" db:"progress_seconds"`
	WatchedAt       time.Time `json:"watched_at" db:"watched_at"`
}
