package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"tech-ip-sem2/services/tasks/client/authclient"
	"tech-ip-sem2/services/tasks/internal/service"
)

type Handler struct {
	svc        *service.TaskService
	authClient *authclient.AuthClient
}

func New(svc *service.TaskService, authClient *authclient.AuthClient) *Handler {
	return &Handler{svc: svc, authClient: authClient}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/v1/tasks", h.tasks)
	mux.HandleFunc("/v1/tasks/", h.taskByID)
}

func (h *Handler) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing token"})
		return false
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	rid := r.Header.Get("X-Request-ID")

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	_, err := h.authClient.Verify(ctx, token, rid)
	if err != nil {
		if strings.Contains(err.Error(), "невалиден") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		} else {
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{"error": "auth unavailable"})
		}
		return false
	}
	return true
}

func (h *Handler) tasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !h.checkAuth(w, r) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(h.svc.List())
	case http.MethodPost:
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			DueDate     string `json:"due_date"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}
		t := h.svc.Create(req.Title, req.Description, req.DueDate)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) taskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !h.checkAuth(w, r) {
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	switch r.Method {
	case http.MethodGet:
		t, ok := h.svc.Get(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
			return
		}
		json.NewEncoder(w).Encode(t)
	case http.MethodPatch:
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Done        *bool  `json:"done"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		t, ok := h.svc.Update(id, req.Title, req.Done, req.Description)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
			return
		}
		json.NewEncoder(w).Encode(t)
	case http.MethodDelete:
		if !h.svc.Delete(id) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
