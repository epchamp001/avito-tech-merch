package dto

import (
	"avito-tech-merch/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapUserInfoResponseToDTO(t *testing.T) {
	now := time.Now()

	purchase1 := &models.Purchase{ID: 1, UserID: 1, MerchID: 101, CreatedAt: now}
	purchase2 := &models.Purchase{ID: 2, UserID: 1, MerchID: 102, CreatedAt: now.Add(time.Minute)}
	transaction1 := &models.Transaction{ID: 10, SenderID: 1, ReceiverID: 2, Amount: 50, CreatedAt: now}
	transaction2 := &models.Transaction{ID: 11, SenderID: 1, ReceiverID: 3, Amount: 75, CreatedAt: now.Add(2 * time.Minute)}

	userInfo := &models.UserInfo{
		UserID:       1,
		Username:     "testuser",
		Balance:      1000,
		Purchases:    []*models.Purchase{purchase1, purchase2},
		Transactions: []*models.Transaction{transaction1, transaction2},
	}

	dto := MapUserInfoResponseToDTO(userInfo)

	assert.Equal(t, userInfo.UserID, dto.UserID)
	assert.Equal(t, userInfo.Username, dto.Username)
	assert.Equal(t, userInfo.Balance, dto.Balance)

	assert.Len(t, dto.Purchases, len(userInfo.Purchases))
	assert.Len(t, dto.Transactions, len(userInfo.Transactions))

	for i, purchase := range userInfo.Purchases {
		expectedDTO := MapPurchaseToDTO(purchase)
		assert.Equal(t, expectedDTO, dto.Purchases[i])
	}

	for i, transaction := range userInfo.Transactions {
		expectedDTO := MapTransactionToDTO(transaction)
		assert.Equal(t, expectedDTO, dto.Transactions[i])
	}
}

func TestMapDTOToUserInfoResponse(t *testing.T) {
	now := time.Now()

	purchaseDTO1 := &PurchaseDTO{ID: 1, UserID: 1, MerchID: 101, CreatedAt: now}
	purchaseDTO2 := &PurchaseDTO{ID: 2, UserID: 1, MerchID: 102, CreatedAt: now.Add(time.Minute)}
	transactionDTO1 := &TransactionDTO{ID: 10, SenderID: 1, ReceiverID: 2, Amount: 50, CreatedAt: now}
	transactionDTO2 := &TransactionDTO{ID: 11, SenderID: 1, ReceiverID: 3, Amount: 75, CreatedAt: now.Add(2 * time.Minute)}

	userInfoDTO := &UserInfoResponse{
		UserID:       1,
		Username:     "testuser",
		Balance:      1000,
		Purchases:    []*PurchaseDTO{purchaseDTO1, purchaseDTO2},
		Transactions: []*TransactionDTO{transactionDTO1, transactionDTO2},
	}

	userInfo := MapDTOToUserInfoResponse(userInfoDTO)

	assert.Equal(t, userInfoDTO.UserID, userInfo.UserID)
	assert.Equal(t, userInfoDTO.Username, userInfo.Username)
	assert.Equal(t, userInfoDTO.Balance, userInfo.Balance)

	assert.Len(t, userInfo.Purchases, len(userInfoDTO.Purchases))
	assert.Len(t, userInfo.Transactions, len(userInfoDTO.Transactions))

	for i, purchaseDTO := range userInfoDTO.Purchases {
		expectedPurchase := MapPurchaseDTOToPurchase(purchaseDTO)
		assert.Equal(t, expectedPurchase, userInfo.Purchases[i])
	}

	for i, transactionDTO := range userInfoDTO.Transactions {
		expectedTransaction := MapTransactionDTOToTransaction(transactionDTO)
		assert.Equal(t, expectedTransaction, userInfo.Transactions[i])
	}
}
