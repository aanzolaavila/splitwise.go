package resources

import "time"

type Group struct {
	Entity
	Name              string        `json:"name"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	Members           []GroupMember `json:"members"`
	SimplifyByDefault bool          `json:"simplify_by_default"`
	OriginalDebts     []Debt        `json:"original_debts"`
	SimplifiedDebts   []Debt        `json:"simplified_debts"`
	Avatar            struct {
		Small    string      `json:"small"`
		Medium   string      `json:"medium"`
		Large    string      `json:"large"`
		Xlarge   string      `json:"xlarge"`
		Xxlarge  string      `json:"xxlarge"`
		Original interface{} `json:"original"`
	} `json:"avatar"`
	TallAvatar struct {
		Xlarge string `json:"xlarge"`
		Large  string `json:"large"`
	} `json:"tall_avatar"`
	CustomAvatar bool `json:"custom_avatar"`
	CoverPhoto   struct {
		Xxlarge string `json:"xxlarge"`
		Xlarge  string `json:"xlarge"`
	} `json:"cover_photo"`
}

type GroupMember struct {
	Entity
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"picture"`
	CustomPicture      bool   `json:"custom_picture"`
	Email              string `json:"email"`
	RegistrationStatus string `json:"registration_status"`
	Balance            []struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"balance"`
}

type Debt struct {
	CurrencyCode string `json:"currency_code"`
	From         int    `json:"from"`
	To           int    `json:"to"`
	Amount       string `json:"amount"`
}
