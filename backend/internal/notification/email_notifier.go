package notification

import (
	"context"
	"log"

	"github.com/sandro/notification-system/internal/domain"
)

type EmailStrategy struct{}

func NewEmailStrategy() *EmailStrategy {
	return &EmailStrategy{}
}

func (s *EmailStrategy) Send(ctx context.Context, user domain.User, message string) error {
	log.Printf("[E-MAIL] -> Sending message to [%s %s]: %s\n", user.Name, user.Email, message)
	log.Printf("[E-MAIL] <- Delivered successfully to %s\n", user.Email)
	return nil
}

func (s *EmailStrategy) GetChannelType() domain.Channel {
	return domain.ChannelEmail
}
