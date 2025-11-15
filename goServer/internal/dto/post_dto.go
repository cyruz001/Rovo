package dto

type CreatePostReq struct {
	Text         string   `json:"text" validate:"required,max=500"`
	ReplyTo      *string  `json:"reply_to" validate:"omitempty,uuid4"`
	IsQuote      bool     `json:"is_quote"`
	QuotedPostID *string  `json:"quoted_post_id" validate:"omitempty,uuid4"`
	MediaURLs    []string `json:"media_urls" validate:"dive,url"`
}

type UpdatePostReq struct {
	Text string `json:"text" validate:"required,max=500"`
}

type PostRes struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	User         UserRes    `json:"user"`
	Text         string     `json:"text"`
	CharCount    int        `json:"char_count"`
	ReplyTo      *string    `json:"reply_to"`
	IsQuote      bool       `json:"is_quote"`
	QuotedPostID *string    `json:"quoted_post_id"`
	Media        []MediaRes `json:"media"`
	LikeCount    int64      `json:"like_count"`
	RepostCount  int64      `json:"repost_count"`
	ReplyCount   int64      `json:"reply_count"`
	IsLiked      bool       `json:"is_liked"`
	IsReposted   bool       `json:"is_reposted"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
}

type PostDetailRes struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	User         UserRes    `json:"user"`
	Text         string     `json:"text"`
	CharCount    int        `json:"char_count"`
	ReplyTo      *string    `json:"reply_to"`
	ReplyToPost  *PostRes   `json:"reply_to_post"`
	IsQuote      bool       `json:"is_quote"`
	QuotedPostID *string    `json:"quoted_post_id"`
	QuotedPost   *PostRes   `json:"quoted_post"`
	Media        []MediaRes `json:"media"`
	LikeCount    int64      `json:"like_count"`
	RepostCount  int64      `json:"repost_count"`
	ReplyCount   int64      `json:"reply_count"`
	IsLiked      bool       `json:"is_liked"`
	IsReposted   bool       `json:"is_reposted"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
}

type MediaRes struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	MediaType string `json:"media_type"`
	Position  int    `json:"position"`
}
