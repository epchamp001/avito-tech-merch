package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
)

type merchService struct {
	repo   db.Repository
	logger logger.Logger
}

func NewMerchService(repo db.Repository, log logger.Logger) MerchService {
	return &merchService{repo: repo, logger: log}
}

func (s *merchService) ListMerch(ctx context.Context) ([]*models.Merch, error) {
	merchList, err := s.repo.GetAllMerch(ctx)
	if err != nil {
		s.logger.Errorw("Failed to fetch merch list",
			"error", err,
		)
		return nil, err
	}

	return merchList, nil
}

func (s *merchService) GetMerch(ctx context.Context, merchID int) (*models.Merch, error) {
	merch, err := s.repo.GetMerchByID(ctx, merchID)
	if err != nil {
		s.logger.Errorw("Failed to fetch merch",
			"merchID", merchID,
			"error", err,
		)
		return nil, err
	}

	return merch, nil
}
