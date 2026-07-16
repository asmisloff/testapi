package state

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// GetToken возвращает токен авторизации из переменной состояния с ключом key
func GetToken(key string) (string, error) {
	if entry, ok := registry[key]; ok {
		return entry.Token, nil
	}
	return "", fmt.Errorf("отсутствует переменная состояния с ключом %s", key)
}

// GetRefreshToken возвращает refresh-токен из переменной состояния с ключом key
func GetRefreshToken(key string) (string, error) {
	if entry, ok := registry[key]; ok {
		return entry.RefreshToken, nil
	}
	return "", fmt.Errorf("отсутствует переменная состояния с ключом %s", key)
}

func init() {
	data, err := os.ReadFile("./state.json")
	if err != nil {
		log.Fatalf("ошибка при попытке чтения файла состояния: %v", err)
	}
	json.Unmarshal(data, &registry)
}

// TestapiState представляет реестр переменных состояния приложения.
// Ключ - строковая метка, значение - запись в реестре.
var registry = make(map[string]*stateEntry)

// stateEntry представляет запись в реестре переменных состояния
type stateEntry struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
