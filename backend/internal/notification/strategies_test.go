package notification

import (
	"context"
	"testing"

	"github.com/sandro/notification-system/internal/domain"
)

func TestEmailStrategy(t *testing.T) {
	strategy := NewEmailStrategy()

	if strategy.GetChannelType() != domain.ChannelEmail {
		t.Errorf("Expected channel type %s, got %s", domain.ChannelEmail, strategy.GetChannelType())
	}

	user := domain.User{Name: "Test User", Email: "test@example.com"}
	err := strategy.Send(context.Background(), user, "Hello Email")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSMSStrategy(t *testing.T) {
	strategy := NewSMSStrategy()

	if strategy.GetChannelType() != domain.ChannelSMS {
		t.Errorf("Expected channel type %s, got %s", domain.ChannelSMS, strategy.GetChannelType())
	}

	user := domain.User{Name: "Test User", PhoneNumber: "123456789"}
	err := strategy.Send(context.Background(), user, "Hello SMS")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestPushStrategy(t *testing.T) {
	strategy := NewPushStrategy()

	if strategy.GetChannelType() != domain.ChannelPush {
		t.Errorf("Expected channel type %s, got %s", domain.ChannelPush, strategy.GetChannelType())
	}

	user := domain.User{Name: "Test User"}
	err := strategy.Send(context.Background(), user, "Hello Push")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
