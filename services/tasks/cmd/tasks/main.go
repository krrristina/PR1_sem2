package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tech-ip-sem2/services/tasks/client/authclient"
	taskshttp "tech-ip-sem2/services/tasks/internal/http"
	"tech-ip-sem2/services/tasks/internal/service"
	"tech-ip-sem2/shared/middleware"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}
	authURL := os.Getenv("AUTH_BASE_URL")
	if authURL == "" {
		authURL = "http://localhost:8081"
	}

	svc := service.New()
	authClient := authclient.New(authURL)
	handler := taskshttp.New(svc, authClient)

	mux := http.NewServeMux()
	handler.Register(mux)

	chain := middleware.RequestID(middleware.Logging("tasks")(mux))

	fmt.Printf("Tasks service запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, chain))
}
