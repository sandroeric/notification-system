package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sandro/notification-system/internal/domain"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	// First fetch all users
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email, phone_number FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	userMap := make(map[int]*domain.User)

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PhoneNumber); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Create pointers map for quick populating
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	// Fetch all subscriptions and populate
	subRows, err := r.db.QueryContext(ctx, "SELECT user_id, category_id FROM user_subscriptions")
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer subRows.Close()

	for subRows.Next() {
		var uid int
		var cat string
		if err := subRows.Scan(&uid, &cat); err != nil {
			return nil, fmt.Errorf("failed to scan sub: %w", err)
		}
		if u, exists := userMap[uid]; exists {
			u.SubscribedCategories = append(u.SubscribedCategories, domain.Category(cat))
		}
	}

	// Fetch all channels and populate
	chRows, err := r.db.QueryContext(ctx, "SELECT user_id, channel_id FROM user_channels")
	if err != nil {
		return nil, fmt.Errorf("failed to query channels: %w", err)
	}
	defer chRows.Close()

	for chRows.Next() {
		var uid int
		var ch string
		if err := chRows.Scan(&uid, &ch); err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		if u, exists := userMap[uid]; exists {
			u.Channels = append(u.Channels, domain.Channel(ch))
		}
	}

	return users, nil
}
