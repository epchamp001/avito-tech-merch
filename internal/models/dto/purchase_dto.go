package dto

import (
	"avito-tech-merch/internal/models"
	"time"
)

// PurchaseDTO DTO for Purchase data in API response
// @Description DTO representing purchase information
type PurchaseDTO struct {
	ID        int       `json:"id" example:"1"`
	UserID    int       `json:"user_id" example:"1"`
	MerchID   int       `json:"merch_id" example:"3"`
	CreatedAt time.Time `json:"created_at" example:"2025-02-15T10:00:00"`
}

// PurchaseSuccessResponse DTO for successful purchase response
// @Description Response indicating that the purchase was successful
type PurchaseSuccessResponse struct {
	Message string `json:"message" example:"purchase successful"`
}

// MapPurchaseToDTO Maps Purchase model to PurchaseDTO
func MapPurchaseToDTO(purchase *models.Purchase) *PurchaseDTO {
	return &PurchaseDTO{
		ID:        purchase.ID,
		UserID:    purchase.UserID,
		MerchID:   purchase.MerchID,
		CreatedAt: purchase.CreatedAt,
	}
}

// MapPurchaseDTOToPurchase Maps PurchaseDTO to Purchase model
func MapPurchaseDTOToPurchase(purchaseDTO *PurchaseDTO) *models.Purchase {
	return &models.Purchase{
		ID:        purchaseDTO.ID,
		UserID:    purchaseDTO.UserID,
		MerchID:   purchaseDTO.MerchID,
		CreatedAt: purchaseDTO.CreatedAt,
	}
}
