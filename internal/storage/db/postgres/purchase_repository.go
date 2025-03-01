package postgres

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresPurchaseRepository struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewPurchaseRepository(pool *pgxpool.Pool, log logger.Logger) db.PurchaseRepository {
	return &postgresPurchaseRepository{pool: pool, logger: log}
}

func (r *postgresPurchaseRepository) CreatePurchase(ctx context.Context, purchase *models.Purchase) (int, error) {
	query := `
        INSERT INTO purchases (user_id, merch_id)
        VALUES ($1, $2)
        RETURNING id
    `

	var purchaseID int
	err := r.pool.QueryRow(ctx, query, purchase.UserID, purchase.MerchID).Scan(&purchaseID)
	if err != nil {
		r.logger.Errorw("Error creating purchase",
			"error", err,
			"userID", purchase.UserID,
			"merchID", purchase.MerchID,
		)
		return 0, fmt.Errorf("failed to create purchase: %w", err)
	}

	return purchaseID, nil
}

func (r *postgresPurchaseRepository) GetPurchaseByUserID(ctx context.Context, userID int) ([]*models.Purchase, error) {
	query := `
        SELECT id, user_id, merch_id, created_at
        FROM purchases
        WHERE user_id = $1
    `

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		r.logger.Errorw("Error retrieving purchase list",
			"error", err,
			"userID", userID,
		)
		return nil, fmt.Errorf("failed to retrieve purchase list: %w", err)
	}
	defer rows.Close()

	var purchases []*models.Purchase
	for rows.Next() {
		var purchase models.Purchase
		err := rows.Scan(
			&purchase.ID,
			&purchase.UserID,
			&purchase.MerchID,
			&purchase.CreatedAt,
		)
		if err != nil {
			r.logger.Errorw("Error scanning purchase data",
				"error", err,
			)
			return nil, fmt.Errorf("error reading purchase data: %w", err)
		}
		purchases = append(purchases, &purchase)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorw("Error processing query result",
			"error", err,
		)
		return nil, fmt.Errorf("error processing query result: %w", err)
	}

	return purchases, nil
}
