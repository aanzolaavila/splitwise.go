package resources

type ExpenseID Identifier

type ExpenseBase struct {
	Cost           string `json:"cost"`
	Description    string `json:"description"`
	Details        string `json:"details"`
	Date           string `json:"date"`
	RepeatInterval string `json:"repeat_interval"`
	CurrencyCode   string `json:"currency_code"`
	CategoryId     uint32 `json:"category_id"`
	GroupId        uint32 `json:"group_id"`
}

type Expense struct {
	ExpenseBase
	ID                     ExpenseID `json:"id"`
	FriendshipID           uint64    `json:"friendship_id"`
	Repeats                bool      `json:"repeats"`
	EmailReminder          bool      `json:"email_reminder"`
	EmailReminderInAdvance int8      `json:"email_reminder_in_advance"`
	NextRepeat             string    `json:"next_repeat"`
	CommentsCount          uint      `json:"comments_count"`
	Payment                bool      `json:"payment"`
	TransactionConfirmed   bool      `json:"transaction_confirmed"`
	CreatedAt              string    `json:"created_at"`
	CreatedBy              User      `json:"created_by"`
	UpdatedAt              string    `json:"updated_at"`
	UpdatedBy              User      `json:"updated_by"`
	DeletedAt              string    `json:"deleted_at"`
	DeletedBy              User      `json:"deleted_by"`
	Repayments             []struct {
		From   uint32 `json:"from"`
		To     uint32 `json:"to"`
		Amount string `json:"amount"`
	} `json:"repayments"`
	Category struct {
		ID   CategoryID `json:"id"`
		Name string     `json:"Name"`
	} `json:"category"`
	Receipt struct {
		Large    string `json:"large"`
		Original string `json:"original"`
	} `json:"receipt"`
	Users []struct {
		User
		UserId     uint64 `json:"user_id"`
		PaidShare  string `json:"paid_share"`
		OwedShare  string `json:"owed_share"`
		NetBalance string `json:"net_balance"`
	} `json:"users"`
	Comments []Comment `json:"comments"`
}
