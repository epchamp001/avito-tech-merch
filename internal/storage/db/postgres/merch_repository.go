package postgres

import (
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresMerchRepository struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewMerchRepository(pool *pgxpool.Pool, log logger.Logger) db.MerchRepository {
	return &postgresMerchRepository{pool: pool, logger: log}
}
