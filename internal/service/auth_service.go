package service

import (
	"avito-tech-merch/internal/config"
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	myjwt "avito-tech-merch/pkg/jwt"
	"avito-tech-merch/pkg/logger"
	pass "avito-tech-merch/pkg/password"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"time"
)

type authService struct {
	repo      db.Repository
	logger    logger.Logger
	JWTConfig config.JWTConfig
}

func NewAuthService(repo db.Repository, log logger.Logger, jwtConfig config.JWTConfig) AuthService {
	return &authService{
		repo:      repo,
		logger:    log,
		JWTConfig: jwtConfig,
	}
}
func (s *authService) Register(ctx context.Context, username string, password string) (string, error) {
	existingUser, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.logger.Errorw("Failed to get user by username",
			"error", err,
			"username", username,
		)
		return "", err
	}
	if existingUser != nil {
		s.logger.Infow("User already exists",
			"username", username,
		)
		return "", fmt.Errorf("user already exists")
	}

	hashedPassword, err := pass.HashPassword(password)
	if err != nil {
		s.logger.Errorw("Failed to hash password",
			"error", err,
			"username", username,
		)
		return "", err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
		Balance:      1000,
		CreatedAt:    time.Now(),
	}

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Errorw("Failed to create user",
			"error", err,
			"username", username,
		)
		return "", err
	}

	token, err := myjwt.GenerateToken(userID, s.JWTConfig.SecretKey, s.JWTConfig.TokenExpiry)
	if err != nil {
		s.logger.Errorw("Failed to generate token",
			"error", err,
			"userID", user.ID,
		)
		return "", err
	}

	return token, nil
}

func (s *authService) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Infow("User not found",
				"username", username,
			)
			return "", errors.New("user not found")
		}
		s.logger.Errorw("Failed to get user by username",
			"error", err,
			"username", username,
		)
		return "", err
	}

	if !pass.CheckPassword(user.PasswordHash, password) {
		s.logger.Infow("Invalid password",
			"username", username,
		)
		return "", fmt.Errorf("invalid password")
	}

	token, err := myjwt.GenerateToken(user.ID, s.JWTConfig.SecretKey, s.JWTConfig.TokenExpiry)
	if err != nil {
		s.logger.Errorw("Failed to generate token",
			"error", err,
			"userID", user.ID,
		)
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *authService) ValidateToken(token string) (int, error) {
	userID, err := myjwt.ParseJWTToken(token, s.JWTConfig.SecretKey)
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			s.logger.Errorw("Invalid token signature",
				"error", err,
			)
			return 0, fmt.Errorf("invalid token signature")
		}
		s.logger.Errorw("Invalid token",
			"error", err,
		)
		return 0, fmt.Errorf("invalid token")
	}

	return userID, nil
}
