//go:build integration

package integration

import (
	"avito-tech-merch/internal/models/dto"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/go-testfixtures/testfixtures/v3"
	"net/http"
)

// TestLoginIntegration_InvalidRequest проверяет, что при некорректном JSON (отсутствует поле password) возвращается HTTP 400
func (s *TestSuite) TestLoginIntegration_InvalidRequest() {
	reqBody := []byte(`{"username": "integration_login_user"}`)
	req, err := http.NewRequest("POST", s.server.URL+"/api/auth/login", bytes.NewBuffer(reqBody))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	var errResp dto.ErrorResponse400
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	s.Require().NoError(err)
	s.Require().Equal(400, errResp.Code)
	s.Require().Equal("invalid request", errResp.Message)
}

// TestLoginIntegration_InvalidCredentials проверяет, что при неверном пароле возвращается HTTP 401
func (s *TestSuite) TestLoginIntegration_InvalidCredentials() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("fixtures/storage"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())

	reqData := dto.LoginRequest{
		Username: "integration_login_user",
		Password: "wrongPassword",
	}
	reqBody, err := json.Marshal(reqData)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.server.URL+"/api/auth/login", bytes.NewBuffer(reqBody))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)

	var errResp dto.ErrorResponseInvalidCredentials401
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	s.Require().NoError(err)
	s.Require().Equal(401, errResp.Code)
	s.Require().Equal("invalid credentials", errResp.Message)
}

// TestLoginIntegration_Success проверяет успешный логин.
func (s *TestSuite) TestLoginIntegration_Success() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("fixtures/storage"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())

	reqData := dto.LoginRequest{
		Username: "integration_login_user",
		Password: "pass123",
	}
	reqBody, err := json.Marshal(reqData)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.server.URL+"/api/auth/login", bytes.NewBuffer(reqBody))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var authResp dto.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	s.Require().NoError(err)
	s.Require().NotEmpty(authResp.Token)
}
