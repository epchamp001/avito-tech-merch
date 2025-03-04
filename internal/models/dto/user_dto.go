package dto

import (
	"avito-tech-merch/internal/models"
)

// UserInfoResponse User information response
// @Description Response containing user information
type UserInfoResponse struct {
	UserID       int               `json:"user_id" example:"1"`
	Username     string            `json:"username" example:"epchamp001"`
	Balance      int               `json:"balance" example:"1500"`
	Purchases    []*PurchaseDTO    `json:"purchases"`
	Transactions []*TransactionDTO `json:"transactions"`
}

// MapUserInfoResponseToDTO Maps internal UserInfo model to UserInfoResponse DTO
func MapUserInfoResponseToDTO(userInfo *models.UserInfo) *UserInfoResponse {
	purchasesDTO := make([]*PurchaseDTO, len(userInfo.Purchases))
	for i, purchase := range userInfo.Purchases {
		purchasesDTO[i] = MapPurchaseToDTO(purchase)
	}

	transactionsDTO := make([]*TransactionDTO, len(userInfo.Transactions))
	for i, transaction := range userInfo.Transactions {
		transactionsDTO[i] = MapTransactionToDTO(transaction)
	}

	return &UserInfoResponse{
		UserID:       userInfo.UserID,
		Username:     userInfo.Username,
		Balance:      userInfo.Balance,
		Purchases:    purchasesDTO,
		Transactions: transactionsDTO,
	}
}

func MapDTOToUserInfoResponse(userInfoDTO *UserInfoResponse) *models.UserInfo {
	purchases := make([]*models.Purchase, len(userInfoDTO.Purchases))
	for i, purchaseDTO := range userInfoDTO.Purchases {
		purchases[i] = MapPurchaseDTOToPurchase(purchaseDTO)
	}

	transactions := make([]*models.Transaction, len(userInfoDTO.Transactions))
	for i, transactionDTO := range userInfoDTO.Transactions {
		transactions[i] = MapTransactionDTOToTransaction(transactionDTO)
	}

	return &models.UserInfo{
		UserID:       userInfoDTO.UserID,
		Username:     userInfoDTO.Username,
		Balance:      userInfoDTO.Balance,
		Purchases:    purchases,
		Transactions: transactions,
	}
}
