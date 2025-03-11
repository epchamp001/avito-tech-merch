package dto

import (
	"avito-tech-merch/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapPurchaseToDTO(t *testing.T) {
	now := time.Now()
	purchase := &models.Purchase{
		ID:        1,
		UserID:    2,
		MerchID:   3,
		CreatedAt: now,
	}

	dto := MapPurchaseToDTO(purchase)

	assert.Equal(t, purchase.ID, dto.ID)
	assert.Equal(t, purchase.UserID, dto.UserID)
	assert.Equal(t, purchase.MerchID, dto.MerchID)
	assert.Equal(t, purchase.CreatedAt, dto.CreatedAt)
}

func TestMapPurchaseDTOToPurchase(t *testing.T) {
	now := time.Now()
	dto := &PurchaseDTO{
		ID:        1,
		UserID:    2,
		MerchID:   3,
		CreatedAt: now,
	}

	purchase := MapPurchaseDTOToPurchase(dto)

	assert.Equal(t, dto.ID, purchase.ID)
	assert.Equal(t, dto.UserID, purchase.UserID)
	assert.Equal(t, dto.MerchID, purchase.MerchID)
	assert.Equal(t, dto.CreatedAt, purchase.CreatedAt)
}
