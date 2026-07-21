// Package auth реализует получение токена доступа для различных типов аутентификации.
// Для типа "sys" используется базовая HTTP-аутентификация с логином и паролем,
// указанными в конфигурации, и возвращается токен из ответа сервера (access_token).
package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"testapi/conf"
	"testapi/state"

	"github.com/golang-jwt/jwt/v5"
)

// GetSystemToken возвращает действующий системный токен для заданного ключа состояния.
// Если срок жизни токена истёк, функция автоматически запрашивает новый токен через
// базовую HTTP-аутентификацию и сохраняет его в состоянии приложения.
func GetSystemToken(key string) (string, error) {
	stateVar, exists := state.Get(key)
	if !exists {
		return "", fmt.Errorf("переменная состояния %q не найдена", key)
	}

	if stateVar.Token != "" && stateVar.ExpiredAt > time.Now().Unix() {
		return stateVar.Token, nil
	}

	newToken, err := requestToken(stateVar.Url, stateVar.Login, stateVar.Password)
	if err != nil {
		return "", fmt.Errorf("ошибка получения нового токена: %w", err)
	}

	exp, err := parseExpiryClaim(newToken)
	if err != nil {
		return "", fmt.Errorf("не удалось определить срок действия нового токена: %w", err)
	}

	stateVar.Token = newToken
	stateVar.ExpiredAt = exp
	if err := state.SaveState(key, stateVar); err != nil {
		return "", fmt.Errorf("ошибка при сохранении файла состояния: %w", err)
	}

	return newToken, nil
}

// GetToken возвращает токен доступа на основе настроек аутентификации.
func GetToken(auth *conf.Auth) (string, error) {
	if auth.Type == "sys" {
		return requestToken(auth.URL, auth.Credentials.Username, auth.Credentials.Password)
	}
	return "", nil
}

// decodeJWTPayload извлекает полезную нагрузку (claims) из JWT‑токена без верификации подписи.
func decodeJWTPayload(tokenStr string) (map[string]any, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("некорректный формат JWT")
	}

	decoded, err := jwt.NewParser().DecodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования payload: %w", err)
	}

	var claims map[string]any
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("ошибка разбора payload: %w", err)
	}

	return claims, nil
}

// parseExpiryClaim извлекает поле exp из JWT‑токена и возвращает его как unix-время.
func parseExpiryClaim(tokenStr string) (int64, error) {
	claims, err := decodeJWTPayload(tokenStr)
	if err != nil {
		return 0, fmt.Errorf("не удалось извлечь payload: %w", err)
	}

	rawExp, ok := claims["exp"]
	if !ok {
		return 0, fmt.Errorf("поле exp отсутствует в payload")
	}

	expFloat, ok := rawExp.(float64)
	if !ok {
		return 0, fmt.Errorf("поле exp имеет некорректный тип %T", rawExp)
	}

	return int64(expFloat), nil
}

// tokenResponse представляет ответ сервера авторизации
type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

var client = http.Client{
	Timeout: 30 * time.Second,
}

// requestToken запрашивает токен у системы авторизации
func requestToken(url string, login string, password string) (string, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(login, password)

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
