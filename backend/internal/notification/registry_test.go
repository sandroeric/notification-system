package notification

import (
	"context"
	"testing"

	"github.com/sandro/notification-system/internal/domain"
)

type dummyStrategy struct {
	channel domain.Channel
}

func (s *dummyStrategy) Send(ctx context.Context, user domain.User, message string) error {
	return nil
}

func (s *dummyStrategy) GetChannelType() domain.Channel {
	return s.channel
}

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	// Initial get should fail
	_, exists := registry.Get(domain.ChannelSMS)
	if exists {
		t.Errorf("Expected SMS strategy to not exist yet")
	}

	// Register a strategy
	smsStrat := &dummyStrategy{channel: domain.ChannelSMS}
	registry.Register(smsStrat)

	// Get should succeed
	retrieved, exists := registry.Get(domain.ChannelSMS)
	if !exists {
		t.Fatalf("Expected SMS strategy to exist")
	}
	if retrieved.GetChannelType() != domain.ChannelSMS {
		t.Errorf("Expected channel %s, got %s", domain.ChannelSMS, retrieved.GetChannelType())
	}
}
