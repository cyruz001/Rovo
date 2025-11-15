package dto

type StatsRes struct {
	TotalUsers    int64 `json:"total_users"`
	TotalPosts    int64 `json:"total_posts"`
	TotalLikes    int64 `json:"total_likes"`
	TotalReposts  int64 `json:"total_reposts"`
	ActiveUsers   int64 `json:"active_users"`
	PostsLastHour int64 `json:"posts_last_hour"`
	PostsLastDay  int64 `json:"posts_last_day"`
	PostsLastWeek int64 `json:"posts_last_week"`
}

type UserStatsRes struct {
	UserID         string  `json:"user_id"`
	TotalPosts     int64   `json:"total_posts"`
	TotalLikes     int64   `json:"total_likes"`
	TotalReposts   int64   `json:"total_reposts"`
	TotalFollowers int64   `json:"total_followers"`
	TotalFollowing int64   `json:"total_following"`
	EngagementRate float64 `json:"engagement_rate"`
}
