package req

import (
	"io"
	"net/http"
)

func Get(url string, params map[string]string) (status int, body string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)

	return resp.StatusCode, string(content), nil
}
