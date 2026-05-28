```go
// Пакет config реализует чтение конфигурационного файла в формате JSON.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// FilterItem описывает один элемент фильтра запроса.
type FilterItem struct {
	Name   string `json:"name"`
	Values []any  `json:"values"`
	Type   string `json:"type"`
}

// Params описывает параметры запроса.
type Params struct {
	Filter []FilterItem   `json:"filter"`
	Extra  map[string]any `json:"-"` // дополнительные поля
}

// UnmarshalJSON реализует кастомный анмаршалинг для Params,
// чтобы обрабатывать произвольные дополнительные поля.
func (p *Params) UnmarshalJSON(data []byte) error {
	// вспомогательная структура для извлечения известных полей
	type knownFields struct {
		Filter []FilterItem `json:"filter"`
	}

	var known knownFields
	if err := json.Unmarshal(data, &known); err != nil {
		return fmt.Errorf("разбор известных полей params: %w", err)
	}
	p.Filter = known.Filter

	// извлекаем все поля как map для дополнительных параметров
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("разбор дополнительных полей params: %w", err)
	}

	delete(raw, "filter")
	p.Extra = raw

	return nil
}

// SetItem описывает один элемент набора запросов.
type SetItem struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Params Params `json:"params"`
}

// Config описывает конфигурацию утилиты.
type Config struct {
	BaseURL string    `json:"base-url"`
	Auth    string    `json:"auth"`
	Set     []SetItem `json:"set"`
}

// Load читает JSON-файл по указанному пути и возвращает структуру Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("чтение файла конфигурации %q: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("разбор файла конфигурации %q: %w", path, err)
	}

	return &cfg, nil
}
```