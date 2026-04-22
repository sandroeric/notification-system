package presentation

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sandro/notification-system/internal/application"
	"github.com/sandro/notification-system/internal/domain"
)

type HTTPHandler struct {
	Service *application.NotificationService
}

func NewHTTPHandler(service *application.NotificationService) *HTTPHandler {
	return &HTTPHandler{Service: service}
}

type NotificationRequest struct {
	Category string `json:"category"`
	Message  string `json:"message"`
}

func (h *HTTPHandler) HandlePostNotification(w http.ResponseWriter, r *http.Request) {
	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.Category = strings.TrimSpace(req.Category)
	req.Message = strings.TrimSpace(req.Message)

	if req.Message == "" {
		http.Error(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}
	if req.Category == "" {
		http.Error(w, "Category cannot be empty", http.StatusBadRequest)
		return
	}

	err := h.Service.Send(r.Context(), domain.Category(req.Category), req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func (h *HTTPHandler) HandleGetLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.Service.GetLogs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// In case logs is nil because it's empty, we explicitly return an empty array
	if logs == nil {
		logs = []domain.NotificationLog{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CorsMiddleware neatly handles preflight OPTIONS requests for our React Vite application
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
