package postgres

import (
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	IsolationLevelSerializable   = pgx.Serializable
	IsolationLevelReadCommitted  = pgx.ReadCommitted
	IsolationLevelRepeatableRead = pgx.RepeatableRead

	AccessModeReadWrite = pgx.ReadWrite
	AccessModeReadOnly  = pgx.ReadOnly
)

type txKey struct{}

type postgresTxManager struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewTxManager(pool *pgxpool.Pool, log logger.Logger) db.TxManager {
	return &postgresTxManager{pool: pool, logger: log}
}

func (t *postgresTxManager) GetExecutor(ctx context.Context) db.Executor {
	if tx, ok := ctx.Value(txKey{}).(db.Executor); ok {
		return tx
	}
	return t.pool
}

func (t *postgresTxManager) WithTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode, fn func(ctx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   isoLevel,
		AccessMode: accessMode,
	}

	var err error

	tx, err := t.pool.BeginTx(ctx, opts)
	if err != nil {
		t.logger.Errorw("Failed to begin transaction",
			"error", err,
		)
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
				t.logger.Errorw("Failed to rollback transaction",
					"error", err,
				)
			}
		}

	}()

	ctx = context.WithValue(ctx, txKey{}, tx)
	if err = fn(ctx); err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		t.logger.Errorw("Failed to commit transaction",
			"error", err,
		)
	}
	return err
}
