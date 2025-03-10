package postgres

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresMerchRepository struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewMerchRepository(pool *pgxpool.Pool, log logger.Logger) db.MerchRepository {
	return &postgresMerchRepository{pool: pool, logger: log}
}

func (r *postgresMerchRepository) GetAllMerch(ctx context.Context) ([]*models.Merch, error) {
	query := `
		SELECT id, name, price 
		FROM merch
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		r.logger.Errorw("Error retrieving merch list",
			"error", err,
		)
		return nil, fmt.Errorf("failed to retrieve merch list: %w", err)
	}
	defer rows.Close()

	var merchList []*models.Merch
	for rows.Next() {
		var merch models.Merch
		err := rows.Scan(
			&merch.ID,
			&merch.Name,
			&merch.Price,
		)
		if err != nil {
			r.logger.Errorw("Error scanning merch data",
				"error", err,
			)
			return nil, fmt.Errorf("error reading merch data: %w", err)
		}
		merchList = append(merchList, &merch)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorw("Error processing query result",
			"error", err,
		)
		return nil, fmt.Errorf("error processing query result: %w", err)
	}

	return merchList, nil
}

func (r *postgresMerchRepository) GetMerchByID(ctx context.Context, id int) (*models.Merch, error) {
	query := `
		SELECT id, name, price 
		FROM merch
		WHERE id = $1
	`

	var merch models.Merch
	err := r.pool.QueryRow(ctx, query, id).Scan(&merch.ID, &merch.Name, &merch.Price)
	if err != nil {
		r.logger.Errorw("Error retrieving merch by ID",
			"error", err,
			"merchID", id,
		)
		return nil, fmt.Errorf("failed to retrieve merch: %w", err)
	}

	return &merch, nil
}

func (r *postgresMerchRepository) GetMerchByName(ctx context.Context, merchName string) (*models.Merch, error) {
	query := `
        SELECT id, name, price
        FROM merch
        WHERE name = $1
    `

	var merch models.Merch
	err := r.pool.QueryRow(ctx, query, merchName).Scan(
		&merch.ID,
		&merch.Name,
		&merch.Price,
	)
	if err != nil {
		r.logger.Errorw("Error retrieving merchandise by name",
			"error", err,
			"merchName", merchName,
		)
		return nil, fmt.Errorf("failed to retrieve merchandise: %w", err)
	}

	return &merch, nil
}
