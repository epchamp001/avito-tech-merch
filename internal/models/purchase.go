package models

import (
	"github.com/google/uuid"
	"time"
)

type Purchase struct {
	ID        uuid.UUID `json:"id"        gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `json:"user_id"   gorm:"type:uuid;not null"`
	MerchID   int       `json:"merch_id"  gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Purchase) TableName() string {
	return "purchases"
}
