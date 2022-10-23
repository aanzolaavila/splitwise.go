package resources

import "time"

type NotificationID Identifier

type Notification struct {
	ID        NotificationID `json:"id"`
	Type      int            `json:"type"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy int            `json:"created_by"`
	Source    struct {
		ID   Identifier `json:"id"`
		Type string     `json:"type"`
		URL  string     `json:"url"`
	} `json:"source"`
	ImageURL   string `json:"image_url"`
	ImageShape string `json:"image_shape"`
	Content    string `json:"content"`
}
