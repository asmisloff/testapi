package req

import (
	"fmt"
	"io"
	"net/http"
	"testapi/auth"
	"testapi/conf"
	"time"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

func ExecRequest(suite *conf.TestSuite, caseId string) error {
	testCase, ok := suite.Cases[caseId]
	if !ok {
		return fmt.Errorf("Тестовый запрос не найден по ID %s\n", caseId)
	}
	token, err := auth.GetToken(&suite.Auth)
	if err != nil {
		return fmt.Errorf("Ошибка при запросе токена авторизации: %w", err)
	}
	switch testCase.Method {
	case "GET":
		status, body, err := Get(suite.BaseURL+testCase.Uri, token, testCase.Params)
		if err != nil {
			return err
		}
		fmt.Printf("%d\n%s\n", status, body)
	default:
		return fmt.Errorf("Неизвестный HTTP метод %s", testCase.Method)
	}
	return nil
}

func Get(url string, token string, params map[string]string) (status int, body string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	return resp.StatusCode, string(content), nil
}
