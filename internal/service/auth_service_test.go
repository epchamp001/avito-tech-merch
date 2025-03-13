package service

import (
	"avito-tech-merch/internal/config"
	"avito-tech-merch/internal/models"
	mockRepo "avito-tech-merch/internal/storage/db/mock"
	mockJWT "avito-tech-merch/pkg/jwt/mock"
	mockLog "avito-tech-merch/pkg/logger/mock"
	pass "avito-tech-merch/pkg/password"
	"context"
	"errors"
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

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"

	txMock := mockRepo.NewTxManager(t)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)
	repoMock.On("GetUserByUsername", ctx, username).Return(nil, pgx.ErrNoRows)
	repoMock.On("CreateUser", ctx, mock.MatchedBy(func(user *models.User) bool {
		return user.Username == username && user.PasswordHash != "" && user.Balance == 1000
	})).Return(1, nil)
	repoMock.On("CommitTx", ctx, txMock).Return(nil)
	expectedToken := "jwt-token"
	jwtMock.On("GenerateToken", 1, jwtConfig.SecretKey, jwtConfig.TokenExpiry).Return(expectedToken, nil)

	token, err := authService.Register(ctx, username, password)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)

	repoMock.AssertCalled(t, "BeginTx", ctx)
	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	repoMock.AssertCalled(t, "CreateUser", ctx, mock.Anything)
	repoMock.AssertCalled(t, "CommitTx", ctx, txMock)
	jwtMock.AssertCalled(t, "GenerateToken", 1, jwtConfig.SecretKey, jwtConfig.TokenExpiry)
}

func TestAuthService_Register_GetUserError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

	ctx := context.Background()
	username := "user"
	password := "pass"
	expectedErr := errors.New("database error")
	txMock := mockRepo.NewTxManager(t)

	repoMock.On("BeginTx", ctx).Return(txMock, nil)
	repoMock.On("GetUserByUsername", ctx, username).Return(nil, expectedErr)
	repoMock.On("RollbackTx", ctx, txMock).Return(nil)
	loggerMock.On("Errorw", "Failed to get user by username",
		"error", expectedErr,
		"username", username,
	).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "BeginTx", ctx)
	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Errorw", "Failed to get user by username",
		"error", expectedErr,
		"username", username,
	)
}

func TestAuthService_Register_UserAlreadyExists(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

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

	txMock := mockRepo.NewTxManager(t)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)
	repoMock.On("GetUserByUsername", ctx, username).Return(existingUser, nil)
	repoMock.On("RollbackTx", ctx, txMock).Return(nil)
	loggerMock.On("Infow", "User already exists", "username", username).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user already exists")
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "BeginTx", ctx)
	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Infow", "User already exists", "username", username)
}

func TestAuthService_Register_CreateUserError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"

	txMock := mockRepo.NewTxManager(t)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)
	repoMock.On("GetUserByUsername", ctx, username).Return(nil, pgx.ErrNoRows)

	expectedErr := errors.New("failed to create user")
	repoMock.On("CreateUser", ctx, mock.MatchedBy(func(user *models.User) bool {
		return user.Username == username && user.PasswordHash != ""
	})).Return(0, expectedErr)
	repoMock.On("RollbackTx", ctx, txMock).Return(nil).Once()

	loggerMock.On("Errorw", "Failed to create user",
		"error", expectedErr,
		"username", username,
	).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "BeginTx", ctx)
	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	repoMock.AssertCalled(t, "CreateUser", ctx, mock.Anything)
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Errorw", "Failed to create user",
		"error", expectedErr,
		"username", username,
	)
}

func TestAuthService_Register_GenerateTokenError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"

	txMock := mockRepo.NewTxManager(t)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)
	repoMock.On("GetUserByUsername", ctx, username).Return(nil, pgx.ErrNoRows)
	repoMock.On("CreateUser", ctx, mock.MatchedBy(func(user *models.User) bool {
		return user.Username == username && user.PasswordHash != ""
	})).Return(1, nil)

	expectedErr := errors.New("failed to generate token")
	jwtMock.On("GenerateToken", 1, jwtConfig.SecretKey, jwtConfig.TokenExpiry).Return("", expectedErr)
	loggerMock.On("Errorw", "Failed to generate token",
		"error", expectedErr,
		"userID", mock.Anything,
	).Return()
	repoMock.On("RollbackTx", ctx, txMock).Return(nil).Once()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "BeginTx", ctx)
	repoMock.AssertCalled(t, "GetUserByUsername", ctx, username)
	repoMock.AssertCalled(t, "CreateUser", ctx, mock.Anything)
	jwtMock.AssertCalled(t, "GenerateToken", 1, jwtConfig.SecretKey, jwtConfig.TokenExpiry)
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Errorw", "Failed to generate token",
		"error", expectedErr,
		"userID", mock.Anything,
	)
}

func TestAuthService_Register_BeginTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"

	expectedErr := errors.New("begin tx error")
	repoMock.On("BeginTx", ctx).Return(nil, expectedErr)
	loggerMock.On("Errorw", "Failed to begin transaction",
		"username", username,
		"error", expectedErr,
	).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "BeginTx", ctx)
	loggerMock.AssertCalled(t, "Errorw", "Failed to begin transaction",
		"username", username,
		"error", expectedErr,
	)
}

func TestAuthService_Register_CommitTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"

	txMock := mockRepo.NewTxManager(t)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)
	repoMock.On("GetUserByUsername", ctx, username).Return(nil, pgx.ErrNoRows)
	repoMock.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(1, nil)
	expectedErr := errors.New("commit tx error")
	repoMock.On("CommitTx", ctx, txMock).Return(expectedErr)
	repoMock.On("RollbackTx", ctx, txMock).Return(nil).Once()
	jwtMock.On("GenerateToken", 1, jwtConfig.SecretKey, jwtConfig.TokenExpiry).Return("jwt-token", nil)
	loggerMock.On("Errorw", "Failed to commit transaction",
		"username", username,
		"error", expectedErr,
	).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to commit transaction")
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "CommitTx", ctx, txMock)
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Errorw", "Failed to commit transaction",
		"username", username,
		"error", expectedErr,
	)
}

func TestAuthService_Register_RollbackTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)
	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	authService := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

	ctx := context.Background()
	username := "newUser"
	password := "password123"

	txMock := mockRepo.NewTxManager(t)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)
	expectedErr := errors.New("get user error")
	repoMock.On("GetUserByUsername", ctx, username).Return(nil, expectedErr)
	rollbackErr := errors.New("rollback failure")
	repoMock.On("RollbackTx", ctx, txMock).Return(rollbackErr).Once()
	loggerMock.On("Errorw", "Failed to get user by username",
		"error", expectedErr,
		"username", username,
	).Return()
	loggerMock.On("Errorw", "Failed to rollback transaction",
		"username", username,
		"error", rollbackErr,
	).Return()

	token, err := authService.Register(ctx, username, password)
	assert.Error(t, err)
	assert.Empty(t, token)

	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Errorw", "Failed to rollback transaction",
		"username", username,
		"error", rollbackErr,
	)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	jwtMock := mockJWT.NewTokenService(t)

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)
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

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)
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

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)
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

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)
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

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)
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

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

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

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

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

	jwtConfig := config.JWTConfig{
		SecretKey:   "secret",
		TokenExpiry: 3600,
	}
	service := NewAuthService(repoMock, loggerMock, jwtConfig, jwtMock)

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
