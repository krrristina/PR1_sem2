package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	authhttp "tech-ip-sem2/services/auth/internal/http"
	"tech-ip-sem2/services/auth/internal/service"
	"tech-ip-sem2/shared/middleware"
)

func main() {
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}

	svc := service.New()
	handler := authhttp.New(svc)

	mux := http.NewServeMux()
	handler.Register(mux)

	chain := middleware.RequestID(middleware.Logging("auth")(mux))

	fmt.Printf("Auth service запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, chain))
}
