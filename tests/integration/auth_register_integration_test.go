package integration

import (
	"avito-tech-merch/internal/models/dto"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/go-testfixtures/testfixtures/v3"
	"net/http"
)

func (s *TestSuite) TestRegisterIntegration_Success() {
	reqData := dto.RegisterRequest{
		Username: "new_integration_user",
		Password: "securePassword123",
	}
	reqBody, err := json.Marshal(reqData)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(reqBody))
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

func (s *TestSuite) TestRegisterIntegration_UserAlreadyExists() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("fixtures/storage"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())

	reqData := dto.RegisterRequest{
		Username: "existing_user",
		Password: "anyPassword",
	}
	reqBody, err := json.Marshal(reqData)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(reqBody))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode)

	var errResp dto.ErrorResponse500
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	s.Require().NoError(err)
	s.Require().Contains(errResp.Message, "user already exists")
}

func (s *TestSuite) TestRegisterIntegration_InvalidRequest() {
	reqBody := []byte(`{"username": "integrationUser"}`)
	req, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(reqBody))
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
