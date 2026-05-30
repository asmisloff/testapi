package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testapi/conf"
	"time"
)

// tokenResponse представляет ответ сервера авторизации
type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

var client = http.Client{
	Timeout: 30 * time.Second,
}

func GetToken(auth *conf.Auth) (string, error) {
	if auth.Type == "sys" {
		return requestToken(auth.URL, &auth.Credentials)
	}
	return "", nil
}

// request запрашивает токен у системы авторизации
func requestToken(url string, creds *conf.AuthCredentials) (string, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(creds.Username, creds.Password)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("Сервер акторизации ответил %d, тело ответа прочитать не удалось: %w", resp.StatusCode, err)
		}
		return "", fmt.Errorf("Сервер авторизации вернул ошибку %d: %s", resp.StatusCode, content)
	}

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("Не удалось декодировать ответ сервера авторизации: %w", err)
	}
	return token.AccessToken, nil
}
