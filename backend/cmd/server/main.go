package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sandro/notification-system/internal/application"
	"github.com/sandro/notification-system/internal/infrastructure"
	"github.com/sandro/notification-system/internal/notification"
	"github.com/sandro/notification-system/internal/presentation"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./notification.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Init Database Repo (executes migrations)
	repo, err := infrastructure.NewSQLiteNotificationRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 2. Setup Notification Registry & Strategies
	registry := notification.NewRegistry()
	registry.Register(notification.NewSMSStrategy())
	registry.Register(notification.NewEmailStrategy())
	registry.Register(notification.NewPushStrategy())

	// 3. Load Mock Users
	users := infrastructure.GetMockUsers()

	// 4. Init Central Service
	svc := application.NewNotificationService(registry, repo, users)

	// 5. Init REST Controllers
	handler := presentation.NewHTTPHandler(svc)

	// 6. Hook up to Router
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/notifications", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.HandlePostNotification(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	
	mux.HandleFunc("/api/notifications/log", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.HandleGetLogs(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Wrap in CORS middleware to permit frontend connectivity
	serverHandler := presentation.CorsMiddleware(mux)

	log.Printf("Backend HTTP Server starting on :%s ...", port)
	if err := http.ListenAndServe(":"+port, serverHandler); err != nil {
		log.Fatalf("Server critically failed: %v", err)
	}
}
