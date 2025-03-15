package service

import (
	"avito-tech-merch/internal/config"
	"avito-tech-merch/internal/models"
	mockRepo "avito-tech-merch/internal/storage/db/mock"
	"avito-tech-merch/internal/storage/db/postgres"
	mockJWT "avito-tech-merch/pkg/jwt/mock"
	mockLog "avito-tech-merch/pkg/logger/mock"
	pass "avito-tech-merch/pkg/password"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestAuthService_Register_Success(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)
	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"
	expectedToken := "jwt-token"

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			// Вызываем транзакционный блок внутри WithTx
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(nil)

	repoMock.
		On("GetUserByUsername", mock.Anything, username).
		Return(nil, pgx.ErrNoRows)
	repoMock.
		On("CreateUser", mock.Anything, mock.MatchedBy(func(user *models.User) bool {
			return user.Username == username && user.PasswordHash != "" && user.Balance == 1000
		})).
		Return(1, nil)
	jwtMock.
		On("GenerateToken", 1, jwtConfig.SecretKey, jwtConfig.TokenExpiry).
		Return(expectedToken, nil)

	token, err := authService.Register(ctx, username, password)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)

	repoMock.AssertCalled(t, "GetUserByUsername", mock.Anything, username)
	repoMock.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything)
	jwtMock.AssertCalled(t, "GenerateToken", 1, jwtConfig.SecretKey, jwtConfig.TokenExpiry)
	txManagerMock.AssertExpectations(t)
}

func TestAuthService_Register_GetUserError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)
	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)

	ctx := context.Background()
	username := "user"
	password := "pass"
	expectedErr := errors.New("database error")

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(expectedErr)

	repoMock.
		On("GetUserByUsername", mock.Anything, username).
		Return(nil, expectedErr)

	loggerMock.
		On("Errorw", "Failed to get user by username",
			"error", expectedErr,
			"username", username,
		).Return()
	loggerMock.
		On("Errorw", "Error during Register operation",
			"error", expectedErr,
		).Return()

	_, err := authService.Register(ctx, username, password)
	assert.Error(t, err)

	repoMock.AssertCalled(t, "GetUserByUsername", mock.Anything, username)
	loggerMock.AssertCalled(t, "Errorw", "Failed to get user by username",
		"error", expectedErr,
		"username", username,
	)
	loggerMock.AssertCalled(t, "Errorw", "Error during Register operation",
		"error", expectedErr,
	)
	txManagerMock.AssertExpectations(t)
}

func TestAuthService_Register_UserAlreadyExists(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)
	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)

	ctx := context.Background()
	username := "existingUser"
	password := "pass"
	existingUser := &models.User{
		ID:           1,
		Username:     username,
		PasswordHash: "hash",
		Balance:      1000,
		CreatedAt:    time.Now(),
	}

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("user already exists"))

	repoMock.
		On("GetUserByUsername", mock.Anything, username).
		Return(existingUser, nil)

	loggerMock.
		On("Infow", "User already exists", "username", username).
		Return()
	loggerMock.
		On("Errorw", "Error during Register operation",
			"error", fmt.Errorf("user already exists"),
		).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user already exists")
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "GetUserByUsername", mock.Anything, username)
	loggerMock.AssertCalled(t, "Infow", "User already exists", "username", username)
	loggerMock.AssertCalled(t, "Errorw", "Error during Register operation",
		"error", fmt.Errorf("user already exists"),
	)
	txManagerMock.AssertExpectations(t)
}

func TestAuthService_Register_CreateUserError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)
	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(errors.New("failed to create user"))

	repoMock.
		On("GetUserByUsername", mock.Anything, username).
		Return(nil, pgx.ErrNoRows)
	repoMock.
		On("CreateUser", mock.Anything, mock.MatchedBy(func(user *models.User) bool {
			return user.Username == username && user.PasswordHash != ""
		})).
		Return(0, errors.New("failed to create user"))

	loggerMock.
		On("Errorw", "Failed to create user",
			"error", errors.New("failed to create user"),
			"username", username,
		).Return()
	loggerMock.
		On("Errorw", "Error during Register operation",
			"error", errors.New("failed to create user"),
		).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "GetUserByUsername", mock.Anything, username)
	repoMock.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything)
	loggerMock.AssertCalled(t, "Errorw", "Failed to create user",
		"error", errors.New("failed to create user"),
		"username", username,
	)
	txManagerMock.AssertExpectations(t)
}

func TestAuthService_Register_GenerateTokenError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)
	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to generate token: %w", errors.New("failed to generate token")))

	repoMock.
		On("GetUserByUsername", mock.Anything, username).
		Return(nil, pgx.ErrNoRows)
	repoMock.
		On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
		Return(1, nil)

	expectedErr := errors.New("failed to generate token")
	jwtMock.
		On("GenerateToken", 1, jwtConfig.SecretKey, jwtConfig.TokenExpiry).
		Return("", expectedErr)

	loggerMock.
		On("Errorw", "Failed to generate token",
			"error", expectedErr,
			"userID", 1,
		).Return()
	loggerMock.
		On("Errorw", "Error during Register operation",
			"error", fmt.Errorf("failed to generate token: %w", expectedErr),
		).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "GetUserByUsername", mock.Anything, username)
	repoMock.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything)
	jwtMock.AssertCalled(t, "GenerateToken", 1, jwtConfig.SecretKey, jwtConfig.TokenExpiry)
	loggerMock.AssertCalled(t, "Errorw", "Failed to generate token",
		"error", expectedErr,
		"userID", 1,
	)
	txManagerMock.AssertExpectations(t)
}

func TestAuthService_Register_BeginTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)

	txManagerMock := mockRepo.NewTxManager(t)
	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"
	expectedErr := errors.New("begin tx error")

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Return(expectedErr)

	loggerMock.
		On("Errorw", "Error during Register operation",
			"error", expectedErr,
		).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Empty(t, token)

	loggerMock.AssertCalled(t, "Errorw", "Error during Register operation",
		"error", expectedErr,
	)
	txManagerMock.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)
	ctx := context.Background()
	username := "nonexistent"
	password := "any-password"

	// Если пользователь не найден, repo возвращает pgx.ErrNoRows
	repoMock.On("GetUserByUsername", ctx, username).Return(nil, pgx.ErrNoRows)
	loggerMock.On("Infow", "User not found", "username", username).Return()

	token, err := service.Login(ctx, username, password)
	assert.Error(t, err)
	assert.EqualError(t, err, "user not found")
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	loggerMock.AssertCalled(t, "Infow", "User not found", "username", username)
}

func TestAuthService_Login_GetUserError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)
	ctx := context.Background()
	username := "testuser"
	password := "password"

	expectedErr := errors.New("db error")
	repoMock.On("GetUserByUsername", ctx, username).Return(nil, expectedErr)
	loggerMock.On("Errorw", "Failed to get user by username", "error", expectedErr, "username", username).Return()

	token, err := service.Login(ctx, username, password)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	loggerMock.AssertCalled(t, "Errorw", "Failed to get user by username", "error", expectedErr, "username", username)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)
	ctx := context.Background()
	username := "testuser"

	correctPassword := "correctpassword"
	wrongPassword := "wrongpassword"

	validHash, err := pass.HashPassword(correctPassword)
	assert.NoError(t, err)

	user := &models.User{
		ID:           1,
		Username:     username,
		PasswordHash: validHash,
		Balance:      100,
	}
	repoMock.On("GetUserByUsername", ctx, username).Return(user, nil)
	loggerMock.On("Infow", "Invalid password", "username", username).Return()

	token, err := service.Login(ctx, username, wrongPassword)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	loggerMock.AssertCalled(t, "Infow", "Invalid password", "username", username)
}

func TestAuthService_Login_TokenError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)
	ctx := context.Background()
	username := "testuser"
	password := "correctpassword"

	validHash, err := pass.HashPassword(password)
	assert.NoError(t, err)
	user := &models.User{
		ID:           1,
		Username:     username,
		PasswordHash: validHash,
		Balance:      100,
	}
	repoMock.On("GetUserByUsername", ctx, username).Return(user, nil)

	expectedErr := errors.New("token error")
	jwtMock.On("GenerateToken", user.ID, jwtConfig.SecretKey, jwtConfig.TokenExpiry).Return("", expectedErr)
	loggerMock.On("Errorw", "Failed to generate token", "error", expectedErr, "userID", user.ID).Return()

	token, err := service.Login(ctx, username, password)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate token")
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	jwtMock.AssertCalled(t, "GenerateToken", user.ID, jwtConfig.SecretKey, jwtConfig.TokenExpiry)
	loggerMock.AssertCalled(t, "Errorw", "Failed to generate token", "error", expectedErr, "userID", user.ID)
}

func TestAuthService_Login_Success(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)
	ctx := context.Background()
	username := "testuser"
	password := "correctpassword"

	validHash, err := pass.HashPassword(password)
	assert.NoError(t, err)
	user := &models.User{
		ID:           1,
		Username:     username,
		PasswordHash: validHash,
		Balance:      100,
	}
	repoMock.On("GetUserByUsername", ctx, username).Return(user, nil)
	expectedToken := "jwt-token"
	jwtMock.On("GenerateToken", user.ID, jwtConfig.SecretKey, jwtConfig.TokenExpiry).Return(expectedToken, nil)

	token, err := service.Login(ctx, username, password)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)

	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	jwtMock.AssertCalled(t, "GenerateToken", user.ID, jwtConfig.SecretKey, jwtConfig.TokenExpiry)
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)

	expectedUserID := 42
	token := "valid-token"

	// Настраиваем токен-сервис: возвращаем корректный userID без ошибки
	jwtMock.On("ParseJWTToken", token, jwtConfig.SecretKey).Return(expectedUserID, nil).Once()

	userID, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)

	jwtMock.AssertCalled(t, "ParseJWTToken", token, jwtConfig.SecretKey)
}

func TestAuthService_ValidateToken_InvalidSignature(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)

	token := "token-with-bad-signature"
	signatureErr := jwt.ErrSignatureInvalid

	// Настраиваем токен-сервис: возвращаем ошибку подписи.
	jwtMock.On("ParseJWTToken", token, jwtConfig.SecretKey).Return(0, signatureErr).Once()
	loggerMock.On("Errorw", "Invalid token signature", "error", signatureErr).Once()

	userID, err := service.ValidateToken(token)
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.EqualError(t, err, "invalid token signature")

	jwtMock.AssertCalled(t, "ParseJWTToken", token, jwtConfig.SecretKey)
	loggerMock.AssertCalled(t, "Errorw", "Invalid token signature", "error", signatureErr)
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	txManagerMock := mockRepo.NewTxManager(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock, txManagerMock)

	token := "invalid-token"
	otherErr := errors.New("some token error")

	// Настраиваем токен-сервис: возвращаем ошибку, отличную от ErrSignatureInvalid.
	jwtMock.On("ParseJWTToken", token, jwtConfig.SecretKey).Return(0, otherErr).Once()
	loggerMock.On("Errorw", "Invalid token", "error", otherErr).Once()

	userID, err := service.ValidateToken(token)
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.EqualError(t, err, "invalid token")

	jwtMock.AssertCalled(t, "ParseJWTToken", token, jwtConfig.SecretKey)
	loggerMock.AssertCalled(t, "Errorw", "Invalid token", "error", otherErr)
}
