package gointrum

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Клиент для запросов к Intrum API
var client = &http.Client{
	Timeout: time.Duration(10 * time.Minute),
}

// Интерфейс структуры API-ответа
type respStruct interface {
	stubInterface()
}

func rawRequest(ctx context.Context, apiKey, u string, p map[string]string, r respStruct) error {
	params := make(url.Values, len(p)+1)
	params.Set("apikey", apiKey) // Параметр, содержащий API-ключ
	for k, v := range p {
		params.Set(k, v)
	}
	httpBody := strings.NewReader(params.Encode()) // Тело запроса

	// Формирование нового запроса

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, httpBody)
	if err != nil {
		return fmt.Errorf("failed to create request for method %s: %w", u, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Отправка запроса на сервер

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request for method %s: %w", u, err)
	}
	defer resp.Body.Close()

	// Обработка ответа от сервера

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from method %s: %w", u, err)
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("status code %d from method %s: %w", resp.StatusCode, u, err)
	}

	// Декодирование ответа

	if err := json.Unmarshal(body, r); err != nil {
		return fmt.Errorf("failed to decode response body from method %s: %w", u, err)
	}

	return nil
}
