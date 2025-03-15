//go:build integration
// +build integration

package integration

import (
	"avito-tech-merch/internal/models/dto"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
)

func (s *TestSuite) TestBuyMerchIntegration_Unauthorized() {
	req, err := http.NewRequest("POST", s.server.URL+"/api/merch/buy/cup", nil)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *TestSuite) TestBuyMerchIntegration_InsufficientFunds() {
	regData := dto.RegisterRequest{
		Username: "buyer_insufficient",
		Password: "password",
	}
	regBody, err := json.Marshal(regData)
	s.Require().NoError(err)

	regReq, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(regBody))
	s.Require().NoError(err)
	regReq.Header.Set("Content-Type", "application/json")

	regResp, err := s.server.Client().Do(regReq)
	s.Require().NoError(err)
	defer regResp.Body.Close()
	s.Require().Equal(http.StatusOK, regResp.StatusCode)

	var authResp dto.AuthResponse
	err = json.NewDecoder(regResp.Body).Decode(&authResp)
	s.Require().NoError(err)
	s.Require().NotEmpty(authResp.Token)

	// Уменьшаем баланс пользователя до значения меньше цены товара
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()
	//  товар "powerbank" стоит 200, а устанавливаем баланс 100
	_, err = db.Exec("UPDATE users SET balance = $1 WHERE username = $2", 100, regData.Username)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.server.URL+"/api/merch/buy/powerbank", nil)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode)

	var errResp dto.ErrorResponse500
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	s.Require().NoError(err)
	s.Require().Contains(errResp.Message, "insufficient funds")
}

func (s *TestSuite) TestBuyMerchIntegration_Success() {
	regData := dto.RegisterRequest{
		Username: "buyer_success",
		Password: "password",
	}
	regBody, err := json.Marshal(regData)
	s.Require().NoError(err)

	regReq, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(regBody))
	s.Require().NoError(err)
	regReq.Header.Set("Content-Type", "application/json")

	regResp, err := s.server.Client().Do(regReq)
	s.Require().NoError(err)
	defer regResp.Body.Close()
	s.Require().Equal(http.StatusOK, regResp.StatusCode)

	var authResp dto.AuthResponse
	err = json.NewDecoder(regResp.Body).Decode(&authResp)
	s.Require().NoError(err)
	s.Require().NotEmpty(authResp.Token)

	req, err := http.NewRequest("POST", s.server.URL+"/api/merch/buy/cup", nil)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var successResp dto.PurchaseSuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&successResp)
	s.Require().NoError(err)
	s.Require().Equal("purchase successful", successResp.Message)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()

	var buyerBalance int
	err = db.QueryRow("SELECT balance FROM users WHERE username=$1", regData.Username).Scan(&buyerBalance)
	s.Require().NoError(err)
	s.Require().Equal(980, buyerBalance)
}

func (s *TestSuite) TestBuyMerchIntegration_InvalidItem() {
	regData := dto.RegisterRequest{
		Username: "buyer_invalid_item",
		Password: "password",
	}
	regBody, err := json.Marshal(regData)
	s.Require().NoError(err)

	regReq, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(regBody))
	s.Require().NoError(err)
	regReq.Header.Set("Content-Type", "application/json")

	regResp, err := s.server.Client().Do(regReq)
	s.Require().NoError(err)
	defer regResp.Body.Close()
	s.Require().Equal(http.StatusOK, regResp.StatusCode)

	var authResp dto.AuthResponse
	err = json.NewDecoder(regResp.Body).Decode(&authResp)
	s.Require().NoError(err)
	s.Require().NotEmpty(authResp.Token)

	// Отправляем запрос на покупку несуществующего товара "nonexistent_item"
	req, err := http.NewRequest("POST", s.server.URL+"/api/merch/buy/nonexistent_item", nil)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	// Ожидаем Internal Server Error, так как товар не найден
	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode)

	var errResp dto.ErrorResponse500
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	s.Require().NoError(err)
	s.Require().Contains(errResp.Message, "failed to get merch")
}
