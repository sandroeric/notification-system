package notification

import (
	"context"
	"log"

	"github.com/sandro/notification-system/internal/domain"
)

type SMSStrategy struct{}

func NewSMSStrategy() *SMSStrategy {
	return &SMSStrategy{}
}

func (s *SMSStrategy) Send(ctx context.Context, user domain.User, message string) error {
	log.Printf("[SMS] -> Attempting to send message to [%s %s]: %s\n", user.Name, user.PhoneNumber, message)
	log.Printf("[SMS] <- Delivered successfully to %s\n", user.PhoneNumber)
	return nil
}

func (s *SMSStrategy) GetChannelType() domain.Channel {
	return domain.ChannelSMS
}
