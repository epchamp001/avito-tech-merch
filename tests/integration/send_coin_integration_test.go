//go:build integration

package integration

import (
	"avito-tech-merch/internal/models/dto"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
)

// TestSendCoin_InvalidRequest проверяет, что если в JSON отсутствуют обязательные поля, возвращается 400
func (s *TestSuite) TestSendCoin_InvalidRequest() {
	regData := dto.RegisterRequest{
		Username: "sender_invalid",
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

	// Отправляем запрос без поля "amount"
	reqBody := []byte(`{"receiver_id": 2}`)
	req, err := http.NewRequest("POST", s.server.URL+"/api/send-coin", bytes.NewBuffer(reqBody))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	var errResp400 dto.ErrorResponse400
	err = json.NewDecoder(resp.Body).Decode(&errResp400)
	s.Require().NoError(err)
	s.Require().Equal(400, errResp400.Code)
	s.Require().Equal("invalid request", errResp400.Message)
}

// TestSendCoin_SameSenderReceiver проверяет, что если отправитель равен получателю, возвращается ошибка
func (s *TestSuite) TestSendCoin_SameSenderReceiver() {
	regData := dto.RegisterRequest{
		Username: "sender_equals_receiver",
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

	// Отправляем запрос, где receiver_id равен id отправителя
	reqData := map[string]interface{}{
		"receiver_id": 1,
		"amount":      100,
	}
	reqBody, err := json.Marshal(reqData)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.server.URL+"/api/send-coin", bytes.NewBuffer(reqBody))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode)

	var errResp500 dto.ErrorResponse500
	err = json.NewDecoder(resp.Body).Decode(&errResp500)
	s.Require().NoError(err)
	s.Require().Contains(errResp500.Message, "cannot transfer to yourself")
}

// TestSendCoin_InsufficientFunds проверяет, что если сумма перевода превышает баланс отправителя, возвращается ошибка.
func (s *TestSuite) TestSendCoin_InsufficientFunds() {
	regData := dto.RegisterRequest{
		Username: "sender_insufficient",
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

	// Попытка перевести сумму, превышающую баланс (например, 1500 > 1000)
	reqData := map[string]interface{}{
		"receiver_id": 2,
		"amount":      1500,
	}
	reqBody, err := json.Marshal(reqData)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.server.URL+"/api/send-coin", bytes.NewBuffer(reqBody))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp.Token)

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode)

	var errResp500 dto.ErrorResponse500
	err = json.NewDecoder(resp.Body).Decode(&errResp500)
	s.Require().NoError(err)
	s.Require().Contains(errResp500.Message, "insufficient funds")
}

// TestSendCoin_Success проверяет успешный перевод монет между двумя пользователями
func (s *TestSuite) TestSendCoin_Success() {
	senderData := dto.RegisterRequest{
		Username: "sender_success",
		Password: "password",
	}
	senderBody, err := json.Marshal(senderData)
	s.Require().NoError(err)
	senderReq, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(senderBody))
	s.Require().NoError(err)
	senderReq.Header.Set("Content-Type", "application/json")
	senderResp, err := s.server.Client().Do(senderReq)
	s.Require().NoError(err)
	defer senderResp.Body.Close()
	s.Require().Equal(http.StatusOK, senderResp.StatusCode)
	var senderAuth dto.AuthResponse
	err = json.NewDecoder(senderResp.Body).Decode(&senderAuth)
	s.Require().NoError(err)
	s.Require().NotEmpty(senderAuth.Token)

	receiverData := dto.RegisterRequest{
		Username: "receiver_success",
		Password: "password",
	}
	receiverBody, err := json.Marshal(receiverData)
	s.Require().NoError(err)
	receiverReq, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(receiverBody))
	s.Require().NoError(err)
	receiverReq.Header.Set("Content-Type", "application/json")
	receiverResp, err := s.server.Client().Do(receiverReq)
	s.Require().NoError(err)
	defer receiverResp.Body.Close()
	s.Require().Equal(http.StatusOK, receiverResp.StatusCode)
	var receiverAuth dto.AuthResponse
	err = json.NewDecoder(receiverResp.Body).Decode(&receiverAuth)
	s.Require().NoError(err)
	s.Require().NotEmpty(receiverAuth.Token)

	reqData := map[string]interface{}{
		"receiver_id": 2,
		"amount":      500,
	}
	reqBody, err := json.Marshal(reqData)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.server.URL+"/api/send-coin", bytes.NewBuffer(reqBody))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+senderAuth.Token)

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var successResp dto.TransferSuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&successResp)
	s.Require().NoError(err)
	s.Require().Equal("coins transferred successfully", successResp.Message)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()

	var senderBalance, receiverBalance int
	err = db.QueryRow("SELECT balance FROM users WHERE username=$1", "sender_success").Scan(&senderBalance)
	s.Require().NoError(err)
	err = db.QueryRow("SELECT balance FROM users WHERE username=$1", "receiver_success").Scan(&receiverBalance)
	s.Require().NoError(err)

	s.Require().Equal(500, senderBalance)
	s.Require().Equal(1500, receiverBalance)
}
