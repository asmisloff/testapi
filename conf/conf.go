// Пакет реализует чтение конфигурационного файла в формате JSON.
package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

// TestSuite представляет набор тестовых запросов
type TestSuite struct {
	BaseURL string          `json:"base-url"`
	Auth    Auth            `json:"auth"`
	Cases   map[string]Case `json:"cases"`
}

// Auth определяет параметры аутентификации
type Auth struct {
	Type        string          `json:"type"`
	URL         string          `json:"url"`
	Credentials AuthCredentials `json:"credentials"`
}

// AuthCredentials содержит логин и пароль для аутентификации
type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Case определяет параметры тестового запроса
type Case struct {
	Method string            `json:"method"`
	Uri    string            `json:"uri"`
	Params map[string]string `json:"params"`
}

// Load читает JSON-файл по указанному пути и возвращает структуру TestSuite.
func Load(path string) (*TestSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("чтение файла %q: %w", path, err)
	}

	testSuite := TestSuite{}
	if err := json.Unmarshal(data, &testSuite); err != nil {
		return nil, fmt.Errorf("разбор файла %q: %w", path, err)
	}

	return &testSuite, nil
}
