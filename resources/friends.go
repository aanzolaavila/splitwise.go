package resources

import "time"

type FriendID Identifier

type Friend struct {
	ID                 FriendID `json:"id"`
	FirstName          string   `json:"first_name"`
	LastName           string   `json:"last_name"`
	Email              string   `json:"email"`
	RegistrationStatus string   `json:"registration_status"`
	Picture            struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"picture"`
	Groups []struct {
		GroupId int `json:"group_id"`
		Balance []struct {
			CurrencyCode string `json:"currency_code"`
			Amount       string `json:"amount"`
		} `json:"balance"`
	} `json:"groups"`
	Balance []struct {
		CurrencyCode string `json:"currency_code"`
		Amount       string `json:"amount"`
	} `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}
