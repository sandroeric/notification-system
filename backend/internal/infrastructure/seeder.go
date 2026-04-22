package infrastructure

import (
	"database/sql"
	"fmt"
)

// SeedDatabase hydrates catalogs and relational users
func SeedDatabase(db *sql.DB) error {
	// 1. Seed Categories
	categories := []string{"Sports", "Finance", "Movies"}
	for _, c := range categories {
		_, err := db.Exec("INSERT OR IGNORE INTO categories (id) VALUES (?)", c)
		if err != nil {
			return fmt.Errorf("failed to seed category %s: %w", c, err)
		}
	}

	// 2. Seed Channels
	channels := []string{"SMS", "E-Mail", "Push Notification"}
	for _, c := range channels {
		_, err := db.Exec("INSERT OR IGNORE INTO channels (id) VALUES (?)", c)
		if err != nil {
			return fmt.Errorf("failed to seed channel %s: %w", c, err)
		}
	}

	// 3. Seed Users
	users := []struct {
		ID    int
		Name  string
		Email string
		Phone string
	}{
		{1, "John Doe", "john@example.com", "+1234567890"},
		{2, "Jane Smith", "jane@example.com", "+0987654321"},
		{3, "Alice Johnson", "alice@example.com", "+1122334455"},
		{4, "Bob Brown", "bob@example.com", "+5544332211"},
	}

	for _, u := range users {
		_, err := db.Exec("INSERT OR IGNORE INTO users (id, name, email, phone_number) VALUES (?, ?, ?, ?)", u.ID, u.Name, u.Email, u.Phone)
		if err != nil {
			return fmt.Errorf("failed to seed user %d: %w", u.ID, err)
		}
	}

	// 4. Seed user_subscriptions
	subs := []struct {
		UserID int
		Cat    string
	}{
		{1, "Sports"}, {1, "Movies"},
		{2, "Finance"},
		{3, "Sports"}, {3, "Finance"}, {3, "Movies"},
		{4, "Movies"},
	}

	for _, s := range subs {
		_, err := db.Exec("INSERT OR IGNORE INTO user_subscriptions (user_id, category_id) VALUES (?, ?)", s.UserID, s.Cat)
		if err != nil {
			return fmt.Errorf("failed to seed sub user %d cat %s: %w", s.UserID, s.Cat, err)
		}
	}

	// 5. Seed user_channels
	chans := []struct {
		UserID int
		Chan   string
	}{
		{1, "SMS"}, {1, "E-Mail"},
		{2, "Push Notification"}, {2, "E-Mail"},
		{3, "SMS"}, {3, "E-Mail"}, {3, "Push Notification"},
		{4, "Push Notification"},
	}

	for _, c := range chans {
		_, err := db.Exec("INSERT OR IGNORE INTO user_channels (user_id, channel_id) VALUES (?, ?)", c.UserID, c.Chan)
		if err != nil {
			return fmt.Errorf("failed to seed channel user %d chan %s: %w", c.UserID, c.Chan, err)
		}
	}

	return nil
}
