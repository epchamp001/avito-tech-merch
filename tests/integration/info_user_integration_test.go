//go:build integration

package integration

import (
	"avito-tech-merch/internal/models/dto"
	"bytes"
	"encoding/json"
	"net/http"
)

// TestGetInfo_Unauthorized_NoToken проверяет, что если запрос отправлен без токена, возвращается 401 с сообщением "missing token".
func (s *TestSuite) TestGetInfo_Unauthorized_NoToken() {
	req, err := http.NewRequest("GET", s.server.URL+"/api/info", nil)
	s.Require().NoError(err)

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)

	var errResp dto.ErrorJWTMissingToken
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	s.Require().NoError(err)
	s.Require().Equal(401, errResp.Code)
	s.Require().Equal("missing token", errResp.Message)
}

// TestGetInfo_Unauthorized_InvalidToken проверяет, что если передан невалидный токен, возвращается 401 с сообщением "invalid token".
func (s *TestSuite) TestGetInfo_Unauthorized_InvalidToken() {
	req, err := http.NewRequest("GET", s.server.URL+"/api/info", nil)
	s.Require().NoError(err)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)

	var errResp dto.ErrorJWTInvalidToken
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	s.Require().NoError(err)
	s.Require().Equal(401, errResp.Code)
	s.Require().Equal("invalid token", errResp.Message)
}

// TestGetInfo_Success проверяет успешное получение информации о пользователе.
func (s *TestSuite) TestGetInfo_Success() {
	reqData := dto.RegisterRequest{
		Username: "integration_info_user",
		Password: "securePassword123",
	}
	reqBody, err := json.Marshal(reqData)
	s.Require().NoError(err)

	regReq, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(reqBody))
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

	infoReq, err := http.NewRequest("GET", s.server.URL+"/api/info", nil)
	s.Require().NoError(err)
	infoReq.Header.Set("Authorization", "Bearer "+authResp.Token)

	infoResp, err := s.server.Client().Do(infoReq)
	s.Require().NoError(err)
	defer infoResp.Body.Close()

	s.Require().Equal(http.StatusOK, infoResp.StatusCode)

	var userInfo dto.UserInfoResponse
	err = json.NewDecoder(infoResp.Body).Decode(&userInfo)
	s.Require().NoError(err)
	s.Require().Equal("integration_info_user", userInfo.Username)
	s.Require().Equal(1, userInfo.UserID)
	s.Require().Equal(1000, userInfo.Balance)
}
