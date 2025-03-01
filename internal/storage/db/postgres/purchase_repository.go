package postgres

import (
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresPurchaseRepository struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewPurchaseRepository(pool *pgxpool.Pool, log logger.Logger) db.PurchaseRepository {
	return &postgresPurchaseRepository{pool: pool, logger: log}
}
