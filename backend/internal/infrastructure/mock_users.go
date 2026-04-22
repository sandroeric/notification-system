package infrastructure

import "github.com/sandro/notification-system/internal/domain"

// GetMockUsers returns a hardcoded array of users serving as an in-memory datastore
func GetMockUsers() []domain.User {
	return []domain.User{
		{
			ID:                   1,
			Name:                 "John Doe",
			Email:                "john@example.com",
			PhoneNumber:          "+1234567890",
			SubscribedCategories: []domain.Category{domain.CategorySports, domain.CategoryMovies},
			Channels:             []domain.Channel{domain.ChannelSMS, domain.ChannelEmail},
		},
		{
			ID:                   2,
			Name:                 "Jane Smith",
			Email:                "jane@example.com",
			PhoneNumber:          "+0987654321",
			SubscribedCategories: []domain.Category{domain.CategoryFinance},
			Channels:             []domain.Channel{domain.ChannelPush, domain.ChannelEmail},
		},
		{
			ID:                   3,
			Name:                 "Alice Johnson",
			Email:                "alice@example.com",
			PhoneNumber:          "+1122334455",
			SubscribedCategories: []domain.Category{domain.CategorySports, domain.CategoryFinance, domain.CategoryMovies},
			Channels:             []domain.Channel{domain.ChannelSMS, domain.ChannelEmail, domain.ChannelPush},
		},
		{
			ID:                   4,
			Name:                 "Bob Brown",
			Email:                "bob@example.com",
			PhoneNumber:          "+5544332211",
			SubscribedCategories: []domain.Category{domain.CategoryMovies},
			Channels:             []domain.Channel{domain.ChannelPush},
		},
	}
}
