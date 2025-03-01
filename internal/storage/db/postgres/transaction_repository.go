package postgres

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresTransactionRepository struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewTransactionRepository(pool *pgxpool.Pool, log logger.Logger) db.TransactionRepository {
	return &postgresTransactionRepository{pool: pool, logger: log}
}

func (r *postgresTransactionRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) (int, error) {
	query := `
        INSERT INTO transactions (sender_id, receiver_id, amount, type)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var transactionID int
	err := r.pool.QueryRow(ctx, query, transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.Type).Scan(&transactionID)
	if err != nil {
		r.logger.Errorw("Error creating transaction",
			"error", err,
			"senderID", transaction.SenderID,
			"receiverID", transaction.ReceiverID,
		)
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transactionID, nil
}

func (r *postgresTransactionRepository) GetTransactionByUserID(ctx context.Context, userID int) ([]*models.Transaction, error) {
	query := `
        SELECT id, sender_id, receiver_id, amount, type, created_at
        FROM transactions
        WHERE sender_id = $1 OR receiver_id = $1
    `

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		r.logger.Errorw("Error retrieving transaction list",
			"error", err,
			"userID", userID,
		)
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
			&transaction.Type,
			&transaction.CreatedAt,
		)
		if err != nil {
			r.logger.Errorw("Error scanning transaction data",
				"error", err,
			)
			return nil, fmt.Errorf("error reading transaction data: %w", err)
		}
		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorw("Error processing query result",
			"error", err,
		)
		return nil, fmt.Errorf("error processing query result: %w", err)
	}

	return transactions, nil
}
