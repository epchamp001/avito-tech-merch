package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage"
	"context"
	"errors"
	"github.com/google/uuid"
)

type merchService struct {
	repo         storage.MerchRepository
	userRepo     storage.UserRepository
	purchaseRepo storage.PurchaseRepository
}

func NewMerchService(repo storage.MerchRepository, userRepo storage.UserRepository, purchaseRepo storage.PurchaseRepository) MerchService {
	return &merchService{repo: repo, userRepo: userRepo, purchaseRepo: purchaseRepo}
}

func (s *merchService) BuyMerch(ctx context.Context, userID uuid.UUID, merchID int) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("пользователь не найден")
	}

	merch, err := s.repo.GetMerchByID(ctx, merchID)
	if err != nil || merch == nil {
		return errors.New("товар не найден")
	}

	if user.Balance < merch.Price {
		return errors.New("недостаточно монет")
	}

	if err := s.userRepo.UpdateBalance(ctx, userID, user.Balance-merch.Price); err != nil {
		return err
	}

	purchase := &models.Purchase{
		UserID:  userID,
		MerchID: merchID,
	}
	return s.purchaseRepo.CreatePurchase(ctx, purchase)
}

func (s *merchService) GetUserPurchases(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error) {
	return s.purchaseRepo.GetPurchasesByUser(ctx, userID)
}

func (s *merchService) GetMerchByName(ctx context.Context, name string) (*models.Merch, error) {
	return s.repo.GetMerchByName(ctx, name)
}
