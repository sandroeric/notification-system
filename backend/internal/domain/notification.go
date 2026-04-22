package domain

import "time"

// DeliveryStatus definitions
type DeliveryStatus string

const (
	DeliveryStatusSuccess DeliveryStatus = "Success"
	DeliveryStatusFailed  DeliveryStatus = "Failed"
)

// NotificationLog represents the delivery record to be saved in the DB
type NotificationLog struct {
	ID             int            `json:"id"`
	Category       Category       `json:"category"`
	Message        string         `json:"message"`
	Channel        Channel        `json:"channel"`
	UserID         int            `json:"user_id"`
	UserEmail      string         `json:"user_email"`
	UserPhone      string         `json:"user_phone"`
	DeliveryStatus DeliveryStatus `json:"delivery_status"`
	ErrorMessage   string         `json:"error_message"` // empty if successful
	RetryCount     int            `json:"retry_count"`
	Timestamp      time.Time      `json:"timestamp"`
}
