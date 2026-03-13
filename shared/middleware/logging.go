package middleware

import (
	"log"
	"net/http"
)

func Logging(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get("X-Request-ID")
			log.Printf("[%s] method=%s path=%s request_id=%s", serviceName, r.Method, r.URL.Path, rid)
			next.ServeHTTP(w, r)
		})
	}
}
