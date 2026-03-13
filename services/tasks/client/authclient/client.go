package authclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AuthClient - клиент для обращения к Auth service
type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

// VerifyResponse - ответ от Auth service
type VerifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject"`
	Error   string `json:"error"`
}

// New создаёт новый клиент
func New(baseURL string) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Second, // ВАЖНО: всегда ставь таймаут!
		},
	}
}

// Verify проверяет токен через Auth service
// Возвращает: subject (кто это), и ошибку (если токен невалиден)
func (c *AuthClient) Verify(ctx context.Context, token string, requestID string) (*VerifyResponse, error) {
	// Создаём запрос с контекстом (для таймаута)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/auth/verify", nil)
	if err != nil {
		return nil, fmt.Errorf("создание запроса: %w", err)
	}

	// Добавляем заголовки
	req.Header.Set("Authorization", "Bearer "+token)
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID) // прокидываем request-id!
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Это может быть таймаут или недоступность Auth
		return nil, fmt.Errorf("запрос к auth: %w", err)
	}
	defer resp.Body.Close()

	// Декодируем ответ
	var verifyResp VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("декодирование ответа: %w", err)
	}

	// Проверяем HTTP-статус
	if resp.StatusCode == http.StatusUnauthorized {
		return &verifyResp, fmt.Errorf("токен невалиден")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth вернул статус %d", resp.StatusCode)
	}

	return &verifyResp, nil
}
