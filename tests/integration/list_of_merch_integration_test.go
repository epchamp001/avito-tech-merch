package integration

//import (
//	"avito-tech-merch/internal/models/dto"
//	"bytes"
//	"encoding/json"
//	"net/http"
//)
//
//// TestListMerchIntegration_Success проверяет, что GET /api/merch возвращает список мерча
//func (s *TestSuite) TestListMerchIntegration_Success() {
//	reqData := dto.RegisterRequest{
//		Username: "user_for_merch",
//		Password: "userPassword",
//	}
//	reqBody, err := json.Marshal(reqData)
//	s.Require().NoError(err)
//
//	regReq, err := http.NewRequest("POST", s.server.URL+"/api/auth/register", bytes.NewBuffer(reqBody))
//	s.Require().NoError(err)
//	regReq.Header.Set("Content-Type", "application/json")
//
//	regResp, err := s.server.Client().Do(regReq)
//	s.Require().NoError(err)
//	defer regResp.Body.Close()
//	s.Require().Equal(http.StatusOK, regResp.StatusCode)
//
//	var authResp dto.AuthResponse
//	err = json.NewDecoder(regResp.Body).Decode(&authResp)
//	s.Require().NoError(err)
//	s.Require().NotEmpty(authResp.Token)
//
//	// Отправляем GET-запрос к /api/merch с валидным токеном.
//	req, err := http.NewRequest("GET", s.server.URL+"/api/merch", nil)
//	s.Require().NoError(err)
//	req.Header.Set("Authorization", "Bearer "+authResp.Token)
//
//	resp, err := s.server.Client().Do(req)
//	s.Require().NoError(err)
//	defer resp.Body.Close()
//
//	s.Require().Equal(http.StatusOK, resp.StatusCode)
//
//	var merchList []*dto.MerchDTO
//	err = json.NewDecoder(resp.Body).Decode(&merchList)
//	s.Require().NoError(err)
//
//	s.Require().Greater(len(merchList), 0, "merch list should not be empty")
//	s.Require().Equal(10, len(merchList))
//	s.Require().Equal("t-shirt", merchList[0].Name)
//	s.Require().Equal(80, merchList[0].Price)
//}
