package resources

import "time"

type User struct {
	Entity
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email"`
	RegistrationStatus string `json:"registration_status"`
	Picture            struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"picture"`
	CustomPicture      bool      `json:"custom_picture"`
	NotificationsRead  time.Time `json:"notifications_read"`
	NotificationsCount int       `json:"notifications_count"`
	Notifications      struct {
		AddedAsFriend bool `json:"added_as_friend"`
	} `json:"notifications"`
	DefaultCurrency string `json:"default_currency"`
	Locale          string `json:"locale"`
	Balance         []struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"balance"`
}
