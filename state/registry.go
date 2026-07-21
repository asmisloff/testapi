// Package state предоставляет реестр для хранения переменных состояния приложения.
// Состояние сохраняется в файле state.json в формате JSON и включает токены авторизации
// и другие параметры, которые могут изменяться во время выполнения.
package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"reflect"
)

// stateFileName определяет имя файла, в который сериализуется реестр состояния.
var stateFileName = "state.json"

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

// SaveState обновляет переменные состояния по ключу key
func SaveState(key string, entry *StateEntry) error {
	registry[key] = entry
	return registry.flush()
}

// StateEntry представляет запись в реестре переменных состояния
type StateEntry struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type stateRegistry map[string]*StateEntry

// registry содержит переменные состояния приложения.
// Ключ: строковая метка, значение: запись в реестре.
var registry = make(stateRegistry)

func init() {
	data, err := os.ReadFile(stateFileName)
	if err != nil {
		if _, match := errors.AsType[*fs.PathError](err); match {
			log.Printf("не найден файл %v. Создан новый.", stateFileName)
			registry.flush()
			return
		}
		log.Fatalf("ошибка при попытке чтения файла состояния: %v, %v", err, reflect.TypeOf(err))
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		log.Fatalf("ошибка разбора содержимого %v: %v", stateFileName, err)
	}
}

// flush сохраняет данные из stateRegistry в файл
func (this *stateRegistry) flush() error {
	json, err := json.Marshal(registry)
	if err != nil {
		return fmt.Errorf("ошибка сериализации реестра переменных состояния: %w", err)
	}
	if err := os.WriteFile(stateFileName, json, 0o644); err != nil {
		return fmt.Errorf("ошибка при сохранении реестра переменных состояния: %w", err)
	}
	return nil
}
