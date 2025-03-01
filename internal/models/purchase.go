package models

import "time"

type Purchase struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	MerchID   int       `json:"merch_id"`
	CreatedAt time.Time `json:"created_at"`
}
