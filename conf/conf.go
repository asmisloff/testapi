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
	Auth    string          `json:"auth"`
	Cases   map[string]Case `json:"cases"`
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
