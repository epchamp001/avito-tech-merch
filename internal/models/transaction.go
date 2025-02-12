package models

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID         uuid.UUID `json:"id"         gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SenderID   uuid.UUID `json:"sender_id"  gorm:"type:uuid"`
	ReceiverID uuid.UUID `json:"receiver_id" gorm:"type:uuid"`
	Amount     int       `json:"amount"     gorm:"check:amount > 0;not null"`
	Type       string    `json:"type"       gorm:"check:type IN ('initial','transfer','purchase');not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
