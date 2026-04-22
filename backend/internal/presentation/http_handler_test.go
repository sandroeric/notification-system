package presentation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sandro/notification-system/internal/domain"
)

type mockNotificationService struct {
	sendFunc    func(ctx context.Context, category domain.Category, message string) error
	getLogsFunc func(ctx context.Context) ([]domain.NotificationLog, error)
}

func (m *mockNotificationService) Send(ctx context.Context, category domain.Category, message string) error {
	if m.sendFunc != nil {
		return m.sendFunc(ctx, category, message)
	}
	return nil
}

func (m *mockNotificationService) GetLogs(ctx context.Context) ([]domain.NotificationLog, error) {
	if m.getLogsFunc != nil {
		return m.getLogsFunc(ctx)
	}
	return nil, nil
}

func TestHandlePostNotification(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		mockSend       func(ctx context.Context, category domain.Category, message string) error
		expectedStatus int
	}{
		{
			name: "Valid Request",
			payload: NotificationRequest{
				Category: string(domain.CategorySports),
				Message:  "Test message",
			},
			mockSend: func(ctx context.Context, category domain.Category, message string) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Category",
			payload: NotificationRequest{
				Category: "Food",
				Message:  "Test message",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty Message",
			payload: NotificationRequest{
				Category: string(domain.CategoryMovies),
				Message:  "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON",
			payload:        "not-json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Partial Success/Errors from Service",
			payload: NotificationRequest{
				Category: string(domain.CategoryFinance),
				Message:  "Stock drop",
			},
			mockSend: func(ctx context.Context, category domain.Category, message string) error {
				return errors.New("partial delivery failed")
			},
			expectedStatus: http.StatusMultiStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockNotificationService{
				sendFunc: tt.mockSend,
			}
			handler := NewHTTPHandler(mockService)

			var body []byte
			if str, ok := tt.payload.(string); ok && str == "not-json" {
				body = []byte(`{"category":`) // malformed JSON
			} else {
				body, _ = json.Marshal(tt.payload)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/notifications", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.HandlePostNotification(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestHandleGetLogs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := &mockNotificationService{
			getLogsFunc: func(ctx context.Context) ([]domain.NotificationLog, error) {
				return []domain.NotificationLog{{ID: 1, Message: "test"}}, nil
			},
		}
		handler := NewHTTPHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/notifications/log", nil)
		w := httptest.NewRecorder()

		handler.HandleGetLogs(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var logs []domain.NotificationLog
		if err := json.NewDecoder(w.Body).Decode(&logs); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(logs) != 1 {
			t.Errorf("expected 1 log, got %d", len(logs))
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := &mockNotificationService{
			getLogsFunc: func(ctx context.Context) ([]domain.NotificationLog, error) {
				return nil, errors.New("db error")
			},
		}
		handler := NewHTTPHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/notifications/log", nil)
		w := httptest.NewRecorder()

		handler.HandleGetLogs(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
