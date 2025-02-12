package service

import "avito-tech-merch/internal/storage"

type service struct {
	UserService
	MerchService
	TransactionService
}

func NewService(userRepo storage.UserRepository, merchRepo storage.MerchRepository, txRepo storage.TransactionRepository, purchaseRepo storage.PurchaseRepository) Service {
	return &service{
		UserService:        NewUserService(userRepo),
		MerchService:       NewMerchService(merchRepo, userRepo, purchaseRepo),
		TransactionService: NewTransactionService(userRepo, txRepo),
	}
}
