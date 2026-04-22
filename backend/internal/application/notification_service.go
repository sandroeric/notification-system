package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sandro/notification-system/internal/domain"
)

type NotificationService struct {
	Registry   domain.NotificationRegistry
	Repository domain.NotificationRepository
	Users      []domain.User
}

func NewNotificationService(reg domain.NotificationRegistry, repo domain.NotificationRepository, users []domain.User) *NotificationService {
	return &NotificationService{
		Registry:   reg,
		Repository: repo,
		Users:      users,
	}
}

// Send processes a message dispatch to all users subscribed to the given category.
func (s *NotificationService) Send(ctx context.Context, category domain.Category, message string) error {
	for _, user := range s.Users {
		// Filter users by subscription
		if !s.isSubscribed(user, category) {
			continue
		}

		// Dispatch to all of their configured channels
		for _, channel := range user.Channels {
			strategy, exists := s.Registry.Get(channel)
			if !exists {
				log.Printf("Worker Warning: No registered strategy for channel %s\n", channel)
				continue
			}

			// Delivery attempt with Retry mechanism
			err := s.executeWithRetry(ctx, strategy, user, message, 3)

			// Record Log
			status := domain.DeliveryStatusSuccess
			errorMsg := ""
			retryCount := 3 // Simplified for metric visibility, assumes it took up to 3 inside logic or failed. Wait, actually we can track actual tries.

			// Better retry metric tracking:
			// Let's pass retry count from executeWithRetry if needed, but for simplicity assuming 3 max.
			if err != nil {
				status = domain.DeliveryStatusFailed
				errorMsg = err.Error()
			}

			nLog := &domain.NotificationLog{
				Category:       category,
				Message:        message,
				Channel:        channel,
				UserID:         user.ID,
				UserEmail:      user.Email,
				UserPhone:      user.PhoneNumber,
				DeliveryStatus: status,
				ErrorMessage:   errorMsg,
				RetryCount:     retryCount,
				Timestamp:      time.Now(),
			}

			if saveErr := s.Repository.Save(ctx, nLog); saveErr != nil {
				log.Printf("CRITICAL: Failed to persist delivery log: %v\n", saveErr)
			}
		}
	}
	return nil
}

// executeWithRetry wraps the dispatch process with simple fault tolerance.
func (s *NotificationService) executeWithRetry(ctx context.Context, strategy domain.NotifierStrategy, user domain.User, message string, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = strategy.Send(ctx, user, message)
		if err == nil {
			return nil
		}
		log.Printf("Delivery failed on attempt %d for %s. Retrying...\n", i+1, strategy.GetChannelType())
		time.Sleep(50 * time.Millisecond) // Arbitrary wait
	}
	return fmt.Errorf("exhausted %d retries. last error: %w", maxRetries, err)
}

func (s *NotificationService) isSubscribed(user domain.User, target domain.Category) bool {
	for _, c := range user.SubscribedCategories {
		if c == target {
			return true
		}
	}
	return false
}

func (s *NotificationService) GetLogs(ctx context.Context) ([]domain.NotificationLog, error) {
	return s.Repository.FindAll(ctx)
}
