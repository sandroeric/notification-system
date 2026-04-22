package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sandro/notification-system/internal/domain"
)

type SQLiteNotificationRepository struct {
	db *sql.DB
}

func NewSQLiteNotificationRepository(dbPath string) (*SQLiteNotificationRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &SQLiteNotificationRepository{db: db}
	if err := repo.migrate(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteNotificationRepository) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS notification_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		channel VARCHAR(50) NOT NULL,
		user_id INTEGER NOT NULL,
		user_email VARCHAR(255),
		user_phone VARCHAR(50),
		delivery_status VARCHAR(50) NOT NULL,
		error_message TEXT,
		retry_count INTEGER NOT NULL,
		timestamp DATETIME NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_timestamp ON notification_logs(timestamp);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteNotificationRepository) Save(ctx context.Context, log *domain.NotificationLog) error {
	query := `
		INSERT INTO notification_logs (
			category, message, channel, user_id, user_email, user_phone, 
			delivery_status, error_message, retry_count, timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, query,
		log.Category,
		log.Message,
		log.Channel,
		log.UserID,
		log.UserEmail,
		log.UserPhone,
		log.DeliveryStatus,
		log.ErrorMessage,
		log.RetryCount,
		log.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to save log: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get insert id: %w", err)
	}
	log.ID = int(id)

	return nil
}

func (r *SQLiteNotificationRepository) FindAll(ctx context.Context) ([]domain.NotificationLog, error) {
	query := `
		SELECT id, category, message, channel, user_id, user_email, user_phone, 
		       delivery_status, error_message, retry_count, timestamp
		FROM notification_logs
		ORDER BY timestamp DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.NotificationLog
	for rows.Next() {
		var log domain.NotificationLog
		err := rows.Scan(
			&log.ID,
			&log.Category,
			&log.Message,
			&log.Channel,
			&log.UserID,
			&log.UserEmail,
			&log.UserPhone,
			&log.DeliveryStatus,
			&log.ErrorMessage,
			&log.RetryCount,
			&log.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
