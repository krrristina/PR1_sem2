package httpx

import (
	"net/http"
	"time"
)

// New создаёт http.Client с обязательным таймаутом
func New(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}
