package model

import "time"

type UserQuest struct {
	UserID     int64      `json:"user_id" db:"user_id"`
	QuestID    int64      `json:"quest_id" db:"quest_id"`
	Progress   int        `json:"progress" db:"progress"`
	CompleteAt *time.Time `json:"complete_at" db:"complete_at"`
}
