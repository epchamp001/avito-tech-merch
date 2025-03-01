package postgres

import (
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresTransactionRepository struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewTransactionRepository(pool *pgxpool.Pool, log logger.Logger) db.TransactionRepository {
	return &postgresTransactionRepository{pool: pool, logger: log}
}
