package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Fetch[T any](c *http.Client, method string, url string, data any, headers map[string]string) (T, error) {
	var result T

	body, err := IntoReader(data)
	if err != nil {
		return result, fmt.Errorf("Fetch: failed parsing into reader: %w", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return result, fmt.Errorf("Fetch: failed creating new request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := c.Do(req)
	if err != nil {
		return result, fmt.Errorf("Fetch: failed sending request: %w", err)
	}
	defer res.Body.Close()

	if !IsSuccessStatusCode(res) {
		resBody, _ := io.ReadAll(res.Body)
		return result, fmt.Errorf("Fetch: unexpected status code: %d: %s", res.StatusCode, string(resBody))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("Fetch: failed reading body: %w", err)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return result, fmt.Errorf("Fetch: failed unmarshalling body: %w", err)
	}

	return result, nil
}

// New function in pkg/fetch/fetch.go
func FormFile[T any](c *http.Client, method, url string, body io.Reader, contentType string, headers map[string]string) (T, error) {
	var result T

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return result, err
	}
	req.Header.Set("Content-Type", contentType)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := c.Do(req)
	if err != nil {
		return result, fmt.Errorf("Fetch: failed sending request: %w", err)
	}
	defer res.Body.Close()

	if !IsSuccessStatusCode(res) {
		resBody, _ := io.ReadAll(res.Body)
		return result, fmt.Errorf("Fetch: unexpected status code: %d: %s", res.StatusCode, string(resBody))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("Fetch: failed reading body: %w", err)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return result, fmt.Errorf("Fetch: failed unmarshalling body: %w", err)
	}

	return result, nil
}

func IntoReader(src any) (*bytes.Reader, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("IntoReader: %w", err)
	}

	return bytes.NewReader(b), nil
}

func IsSuccessStatusCode(res *http.Response) bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}
