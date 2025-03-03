package dto

type TransferRequest struct {
	ReceiverID int `json:"receiver_id" binding:"required"`
	Amount     int `json:"amount" binding:"required"`
}
