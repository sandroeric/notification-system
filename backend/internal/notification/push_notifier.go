package notification

import (
	"context"
	"log"

	"github.com/sandro/notification-system/internal/domain"
)

type PushStrategy struct{}

func NewPushStrategy() *PushStrategy {
	return &PushStrategy{}
}

func (s *PushStrategy) Send(ctx context.Context, user domain.User, message string) error {
	log.Printf("[PUSH] -> Triggering push notification for [%s]: %s\n", user.Name, message)
	log.Printf("[PUSH] <- Push acknowledged by %s\n", user.Name)
	return nil
}

func (s *PushStrategy) GetChannelType() domain.Channel {
	return domain.ChannelPush
}
