package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

// User represents a user account
type User struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Username    string    `gorm:"uniqueIndex;not null" json:"username"`
	DisplayName string    `json:"display_name"`
	Email       string    `gorm:"uniqueIndex;not null" json:"email"`
	Password    string    `gorm:"not null"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	Role        string    `gorm:"default:USER;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`

	// Relations
	Rant          []Rant         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Post          []Post         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Followers     []User         `gorm:"many2many:follows;foreignKey:ID;joinForeignKey:FolloweeID;References:ID;joinReferences:FollowerID"`
	Following     []User         `gorm:"many2many:follows;foreignKey:ID;joinForeignKey:FollowerID;References:ID;joinReferences:FolloweeID"`
	Likes         []Like         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	RePost        []Repost       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Mentions      []Mention      `gorm:"foreignKey:MentionedUserID;constraint:OnDelete:CASCADE"`
	Notifications []Notification `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	RateLimits    []RateLimit    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Rant struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

type Post struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID        string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Text          string    `gorm:"not null" json:"text"`
	CharCount     int       `gorm:"not null" json:"char_count"`
	ReplyTo       *string   `gorm:"type:uuid;index" json:"reply_to"` // null if not a reply
	IsQuote       bool      `gorm:"default:false" json:"is_quote"`
	QuotedTweetID *string   `gorm:"type:uuid;index" json:"quoted_post_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`

	// Relations
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Media       []Media   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Likes       []Like    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Retweets    []Repost  `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Mentions    []Mention `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	RepliedPost *Post     `gorm:"foreignKey:ReplyTo;constraint:OnDelete:CASCADE"`
	QuotedPost  *Post     `gorm:"foreignKey:QuotedTweetID;constraint:OnDelete:CASCADE"`
}

type Media struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	PostID    string    `gorm:"type:uuid;not null;index" json:"post_id"`
	URL       string    `gorm:"not null" json:"url"`
	MediaType string    `json:"media_type"` // image/video
	Position  int       `gorm:"default:0" json:"position"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// Relations
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Like struct {
	UserID    string    `gorm:"type:uuid;not null;primaryKey;index" json:"user_id"`
	PostID    string    `gorm:"type:uuid;not null;primaryKey;index" json:"post_id"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Repost struct {
	UserID    string    `gorm:"type:uuid;not null;primaryKey;index" json:"user_id"`
	PostID    string    `gorm:"type:uuid;not null;primaryKey;index" json:"post_id"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Mention struct {
	ID              string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	PostID          string    `gorm:"type:uuid;not null;index" json:"post_id"`
	MentionedUserID string    `gorm:"type:uuid;not null;index" json:"mentioned_user_id"`
	CreatedAt       time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// Relations
	Post          Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	MentionedUser User `gorm:"foreignKey:MentionedUserID;constraint:OnDelete:CASCADE"`
}

type Notification struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      string    `json:"type"` // e.g., "LIKE", "REPOST", "MENTION"
	Read      bool      `gorm:"default:false" json:"read"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// RateLimit represents rate limiting data
type RateLimit struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Action      string    `gorm:"index" json:"action"`
	WindowStart time.Time `json:"window_start"`
	Count       int       `json:"count"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (t *Post) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (m *Media) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

func (r *RateLimit) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}
