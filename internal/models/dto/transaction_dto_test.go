package dto

import (
	"avito-tech-merch/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapTransactionToDTO(t *testing.T) {
	now := time.Now()
	transaction := &models.Transaction{
		ID:         1,
		SenderID:   10,
		ReceiverID: 20,
		Amount:     150,
		CreatedAt:  now,
	}

	dto := MapTransactionToDTO(transaction)

	assert.Equal(t, transaction.ID, dto.ID)
	assert.Equal(t, transaction.SenderID, dto.SenderID)
	assert.Equal(t, transaction.ReceiverID, dto.ReceiverID)
	assert.Equal(t, transaction.Amount, dto.Amount)
	assert.Equal(t, transaction.CreatedAt, dto.CreatedAt)
}

func TestMapTransactionDTOToTransaction(t *testing.T) {
	now := time.Now()
	dto := &TransactionDTO{
		ID:         1,
		SenderID:   10,
		ReceiverID: 20,
		Amount:     150,
		CreatedAt:  now,
	}

	transaction := MapTransactionDTOToTransaction(dto)

	assert.Equal(t, dto.ID, transaction.ID)
	assert.Equal(t, dto.SenderID, transaction.SenderID)
	assert.Equal(t, dto.ReceiverID, transaction.ReceiverID)
	assert.Equal(t, dto.Amount, transaction.Amount)
	assert.Equal(t, dto.CreatedAt, transaction.CreatedAt)
}
