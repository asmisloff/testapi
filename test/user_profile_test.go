package main

import (
	"testapi/conf"
	"testing"
)

var suiteFile = "./test.json"

func TestGetUserSettings(t *testing.T) {
	suite, err := conf.Load(suiteFile)
	if err != nil {
		t.Fatalf("загрузка конфигурационного файла: %v", err)
	}

	// Проверка базового URL
	if suite.BaseURL != "http://localhost:8080" {
		t.Errorf("ожидался base-url http://localhost:8080, получен %s", suite.BaseURL)
	}

	// Проверка параметров аутентификации
	if suite.Auth.Type != "sys" {
		t.Errorf("ожидался auth.type \"sys\", получен %s", suite.Auth.Type)
	}
	if suite.Auth.URL != "http://localhost:8080/auth/token" {
		t.Errorf("ожидался auth.url http://localhost:8080/auth/token, получен %s", suite.Auth.URL)
	}

	creds := suite.Auth.Credentials
	if creds.Username != "admin" {
		t.Errorf("ожидался username \"admin\", получен %s", creds.Username)
	}
	if creds.Password != "secret" {
		t.Errorf("ожидался password \"secret\", получен %s", creds.Password)
	}

	// Проверка набора тестовых запросов
	if len(suite.Cases) != 1 {
		t.Errorf("ожидалось 1 case, получено %d", len(suite.Cases))
	}

	c01, ok := suite.Cases["01"]
	if !ok {
		t.Errorf("отсутствует case \"01\"")
	} else {
		if c01.Method != "GET" {
			t.Errorf("ожидался метод \"GET\" для case \"01\", получен %s", c01.Method)
		}
		if c01.Uri != "/api/user/settings" {
			t.Errorf("ожидался uri \"/api/user/settings\" для case \"01\", получен %s", c01.Uri)
		}
		if lang, ok := c01.Params["lang"]; !ok || lang != "ru" {
			t.Errorf("ожидался параметр lang \"ru\" для case \"01\", получено %q", lang)
		}
	}
}
