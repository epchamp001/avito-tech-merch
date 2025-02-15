package integration

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func GetJWTToken(t *testing.T) string {
	authData := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}

	payloadBytes, err := json.Marshal(authData)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/auth", bytes.NewReader(payloadBytes))
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	token, ok := result["token"]
	assert.True(t, ok, "Токен не найден в ответе")

	return token
}
