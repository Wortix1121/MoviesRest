package model

import "time"

type Subscriptions struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	PlanID    int64     `json:"plan_id" db:"plan_id"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	IsActive  bool      `json:"is_active" db:"is_active"`
}
