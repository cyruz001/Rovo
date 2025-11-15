package dto

type NotificationRes struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"` // LIKE, REPOST, MENTION, FOLLOW, REPLY
	Message   string `json:"message"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
}

type NotificationDetailRes struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Type      string      `json:"type"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"` // Additional data like post_id, user_id
	Read      bool        `json:"read"`
	CreatedAt string      `json:"created_at"`
}
