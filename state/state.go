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
)

// Get возвращает запись состояния по заданному ключу. Если ключ отсутствует в реестре,
// возвращается нулевая запись и false.
func Get(key string) (StateEntry, bool) {
	if value, ok := registry[key]; ok {
		return *value, true
	}
	return StateEntry{}, false
}

// SaveState обновляет или добавляет запись состояния по ключу key и немедленно сохраняет реестр в файл.
func SaveState(key string, entry StateEntry) error {
	*registry[key] = entry
	return registry.flush()
}

func init() {
	data, err := os.ReadFile(stateFileName)
	if err != nil {
		if _, match := errors.AsType[*fs.PathError](err); match {
			log.Printf("не найден файл %v. Создан новый.", stateFileName)
			registry.flush()
			return
		}
		log.Fatalf("ошибка при попытке чтения файла состояния: %v, %v", err)
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		log.Fatalf("ошибка разбора содержимого %v: %v", stateFileName, err)
	}
}

// flush сохраняет все записи реестра в файл.
func (reg stateRegistry) flush() error {
	json, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации реестра переменных состояния: %w", err)
	}
	if err := os.WriteFile(stateFileName, json, 0o644); err != nil {
		return fmt.Errorf("ошибка при сохранении реестра переменных состояния: %w", err)
	}
	return nil
}
