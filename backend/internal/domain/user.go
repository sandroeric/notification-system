package domain

// Category definitions
type Category string

const (
	CategorySports  Category = "Sports"
	CategoryFinance Category = "Finance"
	CategoryMovies  Category = "Movies"
)

func (c Category) IsValid() bool {
	switch c {
	case CategorySports, CategoryFinance, CategoryMovies:
		return true
	}
	return false
}

// Channel definitions
type Channel string

const (
	ChannelSMS   Channel = "SMS"
	ChannelEmail Channel = "E-Mail"
	ChannelPush  Channel = "Push Notification"
)

// User represents the domain model for a mocked user
type User struct {
	ID                   int
	Name                 string
	Email                string
	PhoneNumber          string
	SubscribedCategories []Category
	Channels             []Channel
}
