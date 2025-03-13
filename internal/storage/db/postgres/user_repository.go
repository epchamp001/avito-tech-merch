package postgres

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type postgresUserRepository struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewUserRepository(pool *pgxpool.Pool, log logger.Logger) db.UserRepository {
	return &postgresUserRepository{pool: pool, logger: log}
}

func (r *postgresUserRepository) CreateUser(ctx context.Context, user *models.User) (int, error) {
	query := `
		INSERT INTO users (username, password_hash, balance, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	user.CreatedAt = time.Now()

	var userID int
	err := r.pool.QueryRow(ctx, query, user.Username, user.PasswordHash, user.Balance, user.CreatedAt).Scan(&userID)
	if err != nil {
		r.logger.Errorw("Error when creating a user",
			"error", err,
			"username", user.Username,
		)
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (r *postgresUserRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, balance, created_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
	)
	if err != nil {
		r.logger.Errorw("Error when getting a user by ID",
			"error", err,
			"userID", userID,
		)
		return nil, fmt.Errorf("failed to get a user by ID: %w", err)
	}

	return &user, nil
}

func (r *postgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, balance, created_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
	)
	if err != nil {
		r.logger.Warnw("Error when getting a user by username",
			"error", err,
			"username", username,
		)
		return nil, fmt.Errorf("failed to get a user by username: %w", err)
	}

	return &user, nil
}

func (r *postgresUserRepository) GetBalanceByID(ctx context.Context, userID int) (int, error) {
	query := `
		SELECT balance
		FROM users
		WHERE id = $1
	`

	var balance int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&balance)
	if err != nil {
		r.logger.Errorw("Error when getting a user balance by userID",
			"error", err,
			"userID", userID)
		return 0, fmt.Errorf("failed to get a user balance by userID: %w", err)
	}

	return balance, nil
}

func (r *postgresUserRepository) GetBalanceByName(ctx context.Context, username string) (int, error) {
	query := `
		SELECT balance
		FROM users
		WHERE id = $1
	`

	var balance int
	err := r.pool.QueryRow(ctx, query, username).Scan(&balance)
	if err != nil {
		r.logger.Errorw("Error when getting a user balance by username",
			"error", err,
			"username", username)
		return 0, fmt.Errorf("failed to get a user balance by username: %w", err)
	}

	return balance, nil
}

func (r *postgresUserRepository) UpdateBalance(ctx context.Context, userID int, newBalance int) error {
	query := `
		UPDATE users
		SET balance = $1
		WHERE id = $2
	`

	result, err := r.pool.Exec(ctx, query, newBalance, userID)
	if err != nil {
		r.logger.Errorw("Error when updating a user balance",
			"error", err,
			"userID", userID,
			"newBalance", newBalance,
		)
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warnw("User not found",
			"userID", userID,
		)
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}

func (r *postgresUserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, username, password_hash, balance, created_at
		FROM users
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		r.logger.Errorw("Error retrieving user list",
			"error", err,
		)
		return nil, fmt.Errorf("failed to retrieve user list: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.PasswordHash,
			&user.Balance,
			&user.CreatedAt,
		)
		if err != nil {
			r.logger.Errorw("Error scanning user data",
				"error", err,
			)
			return nil, fmt.Errorf("error reading user data: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorw("Error processing query result",
			"error", err,
		)
		return nil, fmt.Errorf("error processing query result: %w", err)
	}

	return users, nil
}

func (r *postgresUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, password_hash = $2, balance = $3
		WHERE id = $4
	`

	result, err := r.pool.Exec(ctx, query, user.Username, user.PasswordHash, user.Balance, user.ID)
	if err != nil {
		r.logger.Errorw("Error updating user data",
			"error", err,
			"userID", user.ID,
		)
		return fmt.Errorf("failed to update user data: %w", err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warnw("User not found",
			"userID", user.ID,
		)
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	return nil
}
