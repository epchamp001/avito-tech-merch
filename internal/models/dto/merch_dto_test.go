package dto

import (
	"avito-tech-merch/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapMerchToDTO(t *testing.T) {
	merch := &models.Merch{
		ID:    1,
		Name:  "Test Merch",
		Price: 100,
	}

	dto := MapMerchToDTO(merch)

	assert.Equal(t, merch.ID, dto.ID)
	assert.Equal(t, merch.Name, dto.Name)
	assert.Equal(t, merch.Price, dto.Price)
}

func TestMapMerchDTOToMerch(t *testing.T) {
	dto := &MerchDTO{
		ID:    1,
		Name:  "Test Merch",
		Price: 100,
	}

	merch := MapMerchDTOToMerch(dto)

	assert.Equal(t, dto.ID, merch.ID)
	assert.Equal(t, dto.Name, merch.Name)
	assert.Equal(t, dto.Price, merch.Price)
}
