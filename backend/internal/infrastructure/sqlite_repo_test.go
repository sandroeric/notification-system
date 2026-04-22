package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/sandro/notification-system/internal/domain"
)

func TestSQLiteNotificationRepository(t *testing.T) {
	repo, err := NewSQLiteNotificationRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to init in-memory db: %v", err)
	}

	ctx := context.Background()
	timestamp := time.Now()

	logEntry := &domain.NotificationLog{
		Category:       domain.CategorySports,
		Message:        "Test Message",
		Channel:        domain.ChannelEmail,
		UserID:         1,
		UserEmail:      "test@example.com",
		UserPhone:      "123",
		DeliveryStatus: domain.DeliveryStatusSuccess,
		ErrorMessage:   "",
		RetryCount:     0,
		Timestamp:      timestamp,
	}

	// Test Save
	err = repo.Save(ctx, logEntry)
	if err != nil {
		t.Fatalf("failed to save log: %v", err)
	}

	if logEntry.ID == 0 {
		t.Errorf("expected ID to be set after save")
	}

	// Test FindAll
	logs, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("failed to find logs: %v", err)
	}

	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	savedLog := logs[0]
	if savedLog.Message != "Test Message" {
		t.Errorf("expected msg 'Test Message', got '%s'", savedLog.Message)
	}
	if savedLog.Category != domain.CategorySports {
		t.Errorf("expected category 'Sports', got '%s'", savedLog.Category)
	}
}
