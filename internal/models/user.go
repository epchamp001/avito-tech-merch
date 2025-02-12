package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"password" gorm:"not null"` // В JSON лучше не возвращать пароль, поэтому можно делать json:"-"
	Balance   int       `json:"balance"  gorm:"not null;default:1000"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
