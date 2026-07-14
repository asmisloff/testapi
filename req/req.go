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
		return fmt.Errorf("тестовый запрос не найден по ID %s", caseId)
	}
	token, err := auth.GetToken(&suite.Auth)
	if err != nil {
		return fmt.Errorf("ошибка при запросе токена авторизации: %w", err)
	}
	switch testCase.Method {
	case "GET":
		status, body, err := Get(suite.BaseURL+testCase.Uri, token, testCase.Params)
		if err != nil {
			return fmt.Errorf("ошибка при выполнении GET запроса: %w", err)
		}
		fmt.Printf("%d\n%s\n", status, body)
	default:
		return fmt.Errorf("неизвестный HTTP метод %s", testCase.Method)
	}
	return nil
}

// Get выполняет GET запрос с авторизацией и опциональными параметрами.
func Get(url string, token string, params map[string]string) (status int, body string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	// Добавляем query-параметры
	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, "", fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	return resp.StatusCode, string(content), nil
}
