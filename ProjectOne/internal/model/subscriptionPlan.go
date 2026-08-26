package model

import "encoding/json"

type SubscriptionPlan struct {
	ID       int64           `json:"id" db:"id"`
	Name     string          `json:"name" db:"name"`
	Price    float64         `json:"price" db:"price"`
	Features json.RawMessage `json:"features" db:"features"`
}
