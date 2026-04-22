package application

import (
	"context"
	"errors"
	"testing"

	"github.com/sandro/notification-system/internal/domain"
	"github.com/sandro/notification-system/internal/notification"
)

// MockStrategy
type MockStrategy struct {
	Channel   domain.Channel
	FailCount int // how many times to fail before succeeding
	calls     int
}

func (m *MockStrategy) Send(ctx context.Context, user domain.User, message string) error {
	m.calls++
	if m.calls <= m.FailCount {
		return errors.New("simulated failure")
	}
	return nil
}

func (m *MockStrategy) GetChannelType() domain.Channel {
	return m.Channel
}

// MockRepo
type MockRepo struct {
	Logs []domain.NotificationLog
}

func (m *MockRepo) Save(ctx context.Context, log *domain.NotificationLog) error {
	m.Logs = append(m.Logs, *log)
	return nil
}

func (m *MockRepo) FindAll(ctx context.Context) ([]domain.NotificationLog, error) {
	return m.Logs, nil
}

func TestNotificationService_Send(t *testing.T) {
	ctx := context.Background()

	users := []domain.User{
		{
			ID:                   1,
			Email:                "t1@test.com",
			SubscribedCategories: []domain.Category{domain.CategorySports},
			Channels:             []domain.Channel{domain.ChannelEmail},
		},
	}

	repo := &MockRepo{}
	reg := notification.NewRegistry()
	
	// Perfect Strategy
	emailStrat := &MockStrategy{Channel: domain.ChannelEmail, FailCount: 0}
	reg.Register(emailStrat)

	svc := NewNotificationService(reg, repo, users)

	err := svc.Send(ctx, domain.CategorySports, "Test message!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.Logs) != 1 {
		t.Fatalf("expected 1 log persisted, got %d", len(repo.Logs))
	}

	if emailStrat.calls != 1 {
		t.Fatalf("expected 1 call to email strategy, got %d", emailStrat.calls)
	}
}

func TestNotificationService_Send_WithRetries(t *testing.T) {
	ctx := context.Background()

	users := []domain.User{
		{
			ID:                   2,
			Channels:             []domain.Channel{domain.ChannelSMS},
			SubscribedCategories: []domain.Category{domain.CategoryMovies},
		},
	}

	repo := &MockRepo{}
	reg := notification.NewRegistry()
	
	smsStrat := &MockStrategy{Channel: domain.ChannelSMS, FailCount: 2} // Fails twice, succeeds 3rd
	reg.Register(smsStrat)

	svc := NewNotificationService(reg, repo, users)

	svc.Send(ctx, domain.CategoryMovies, "Movie Time")

	if smsStrat.calls != 3 {
		t.Fatalf("expected strategy to be called 3 times due to retries, got %d", smsStrat.calls)
	}
	if repo.Logs[0].DeliveryStatus != domain.DeliveryStatusSuccess {
		t.Fatalf("expected successful delivery status after retry recovery")
	}
}

func TestNotificationService_Send_Failed(t *testing.T) {
	ctx := context.Background()

	users := []domain.User{
		{
			ID:                   3,
			Channels:             []domain.Channel{domain.ChannelPush},
			SubscribedCategories: []domain.Category{domain.CategoryFinance},
		},
	}

	repo := &MockRepo{}
	reg := notification.NewRegistry()
	
	pushStrat := &MockStrategy{Channel: domain.ChannelPush, FailCount: 5} // Will exhaust 3 retries
	reg.Register(pushStrat)

	svc := NewNotificationService(reg, repo, users)
	svc.Send(ctx, domain.CategoryFinance, "Market Crash")

	if pushStrat.calls != 3 {
		t.Fatalf("expected 3 max calls, got %d", pushStrat.calls)
	}
	
	log := repo.Logs[0]
	if log.DeliveryStatus != domain.DeliveryStatusFailed {
		t.Fatalf("expected failure status, got %v", log.DeliveryStatus)
	}
	if log.ErrorMessage == "" {
		t.Fatalf("expected an error message to be captured")
	}
}
