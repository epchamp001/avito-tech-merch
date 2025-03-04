package dto

import "avito-tech-merch/internal/models"

// MerchDTO DTO for Merch data in API response
// @Description DTO representing merch information
type MerchDTO struct {
	ID    int    `json:"id" example:"2"`
	Name  string `json:"name" example:"cup"`
	Price int    `json:"price" example:"20"`
}

// MapMerchToDTO Maps Merch model to MerchDTO
func MapMerchToDTO(merch *models.Merch) *MerchDTO {
	return &MerchDTO{
		ID:    merch.ID,
		Name:  merch.Name,
		Price: merch.Price,
	}
}

// MapMerchDTOToMerch Maps MerchDTO to Merch model
func MapMerchDTOToMerch(merchDTO *MerchDTO) *models.Merch {
	return &models.Merch{
		ID:    merchDTO.ID,
		Name:  merchDTO.Name,
		Price: merchDTO.Price,
	}
}
