package domain

import "context"

// NotifierStrategy handles the actual sending of notifications to a specific channel
type NotifierStrategy interface {
	Send(ctx context.Context, user User, message string) error
	GetChannelType() Channel
}

// NotificationRegistry manages available NotifierStrategies and retrieves them
type NotificationRegistry interface {
	Register(strategy NotifierStrategy)
	Get(channel Channel) (NotifierStrategy, bool)
}

// NotificationRepository handles persistence of notification delivery logs
type NotificationRepository interface {
	Save(ctx context.Context, log *NotificationLog) error
	FindAll(ctx context.Context) ([]NotificationLog, error)
}
