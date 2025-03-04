package dto

import (
	"avito-tech-merch/internal/models"
	"time"
)

// TransactionDTO DTO for Transaction data in API response
// @Description DTO representing transaction information
type TransactionDTO struct {
	ID         int       `json:"id" example:"1"`
	SenderID   int       `json:"sender_id" example:"1"`
	ReceiverID int       `json:"receiver_id" example:"2"`
	Amount     int       `json:"amount" example:"200"`
	CreatedAt  time.Time `json:"created_at" example:"2025-02-16T14:30:00"`
}

// TransferSuccessResponse DTO for successful coin transfer response
// @Description Response indicating that the coin transfer was successful
type TransferSuccessResponse struct {
	Message string `json:"message" example:"coins transferred successfully"`
}

// MapTransactionToDTO Maps Transaction model to TransactionDTO
func MapTransactionToDTO(transaction *models.Transaction) *TransactionDTO {
	return &TransactionDTO{
		ID:         transaction.ID,
		SenderID:   transaction.SenderID,
		ReceiverID: transaction.ReceiverID,
		Amount:     transaction.Amount,
		CreatedAt:  transaction.CreatedAt,
	}
}

// MapTransactionDTOToTransaction Maps TransactionDTO to Transaction model
func MapTransactionDTOToTransaction(transactionDTO *TransactionDTO) *models.Transaction {
	return &models.Transaction{
		ID:         transactionDTO.ID,
		SenderID:   transactionDTO.SenderID,
		ReceiverID: transactionDTO.ReceiverID,
		Amount:     transactionDTO.Amount,
		CreatedAt:  transactionDTO.CreatedAt,
	}
}
