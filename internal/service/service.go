package service

import "avito-tech-merch/internal/storage"

type service struct {
	UserService
	MerchService
	TransactionService
}

func NewService(repo storage.Repository) Service {
	return &service{
		UserService:        NewUserService(repo),
		MerchService:       NewMerchService(repo, repo, repo),
		TransactionService: NewTransactionService(repo, repo),
	}
}
