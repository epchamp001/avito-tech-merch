package postgres

import (
	"avito-tech-merch/internal/metrics"
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"time"
)

type postgresTransactionRepository struct {
	conn   db.TxManager
	logger logger.Logger
}

func NewTransactionRepository(conn db.TxManager, log logger.Logger) db.TransactionRepository {
	return &postgresTransactionRepository{conn: conn, logger: log}
}

func (r *postgresTransactionRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) (int, error) {
	start := time.Now()
	defer func() {
		metrics.RecordDBQueryDuration("CreateTransaction", time.Since(start).Seconds())
	}()

	pool := r.conn.GetExecutor(ctx)

	query := `
        INSERT INTO transactions (sender_id, receiver_id, amount, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var transactionID int
	err := pool.QueryRow(ctx, query, transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.CreatedAt).Scan(&transactionID)
	if err != nil {
		r.logger.Errorw("Error creating transaction",
			"error", err,
			"senderID", transaction.SenderID,
			"receiverID", transaction.ReceiverID,
		)
		metrics.RecordDBError("CreateTransaction")
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transactionID, nil
}

func (r *postgresTransactionRepository) GetTransactionByUserID(ctx context.Context, userID int) ([]*models.Transaction, error) {
	start := time.Now()
	defer func() {
		metrics.RecordDBQueryDuration("GetTransactionByUserID", time.Since(start).Seconds())
	}()

	pool := r.conn.GetExecutor(ctx)

	query := `
        SELECT id, sender_id, receiver_id, amount, created_at
        FROM transactions
        WHERE sender_id = $1 OR receiver_id = $1
    `

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		r.logger.Errorw("Error retrieving transaction list",
			"error", err,
			"userID", userID,
		)
		metrics.RecordDBError("GetTransactionByUserID")
		return nil, fmt.Errorf("failed to retrieve transaction list: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.SenderID,
			&transaction.ReceiverID,
			&transaction.Amount,
			&transaction.CreatedAt,
		)
		if err != nil {
			r.logger.Errorw("Error scanning transaction data",
				"error", err,
			)
			metrics.RecordDBError("GetTransactionByUserID")
			return nil, fmt.Errorf("error reading transaction data: %w", err)
		}
		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorw("Error processing query result",
			"error", err,
		)
		metrics.RecordDBError("GetTransactionByUserID")
		return nil, fmt.Errorf("error processing query result: %w", err)
	}

	return transactions, nil
}
