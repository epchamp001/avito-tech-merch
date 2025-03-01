package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Balance      int       `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
}
