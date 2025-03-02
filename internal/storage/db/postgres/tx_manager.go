package postgres

import (
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresTxManager struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewTxManager(pool *pgxpool.Pool, log logger.Logger) db.TxManager {
	return &postgresTxManager{pool: pool, logger: log}
}

func (p *postgresTxManager) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

func (p *postgresTxManager) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (p *postgresTxManager) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}
