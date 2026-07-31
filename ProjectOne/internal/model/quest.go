package model

type Quest struct {
	ID            int64  `json:"id" db:"id"`
	Title         string `json:"title" db:"title"`
	Description   string `json:"description" db:"description"`
	Type          string `json:"type" db:"type"`
	RequiredCount int64  `json:"required_count" db:"required_count"`
	RewardPoints  int64  `json:"reward_points" db:"reward_points"`
}
