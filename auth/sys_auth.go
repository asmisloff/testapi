// Package auth реализует получение токена доступа для различных типов аутентификации.
// Для типа "sys" используется базовая HTTP-аутентификация с логином и паролем,
// указанными в конфигурации, и возвращается токен из ответа сервера (access_token).
package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testapi/conf"
	"time"
)

const authFilePath = "auth"

// GetToken возвращает токен доступа на основе настроек аутентификации.
// Если в конфигурации указан тип "sys", запрос отправляется с Basic‑аутентификацией,
// используя логин и пароль из параметров Credentials, и возвращается access_token из ответа.
// Для других типов возвращается пустая строка и nil-ошибка.
func GetToken(auth *conf.Auth) (string, error) {
	if auth.Type == "sys" {
		// content, err := os.ReadFile(authFilePath)
		// if err != nil {
		// 	return "", fmt.Errorf("чтение файла %q: %w", authFilePath, err)
		// }
		// authRegistry := make(internal.AuthRegistry)
		// json.Unmarshal(content, &authRegistry)
		// if val, ok := authRegistry[auth.URL]; ok {
		// 	return val.Token, nil
		// }
		return requestToken(auth.URL, &auth.Credentials)
	}
	return "", nil
}

// tokenResponse представляет ответ сервера авторизации
type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

var client = http.Client{
	Timeout: 30 * time.Second,
}

// requestToken запрашивает токен у системы авторизации
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
			return "", fmt.Errorf("Сервер авторизации ответил %d, тело ответа прочитать не удалось: %w", resp.StatusCode, err)
		}
		return "", fmt.Errorf("Сервер авторизации вернул ошибку %d: %s", resp.StatusCode, content)
	}

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("Не удалось декодировать ответ сервера авторизации: %w", err)
	}
	return token.AccessToken, nil
}
