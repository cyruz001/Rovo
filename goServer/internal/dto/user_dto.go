package dto

type UserRegisterReq struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateReq struct {
	DisplayName string `json:"display_name" validate:"max=50"`
	Bio         string `json:"bio" validate:"max=500"`
	AvatarURL   string `json:"avatar_url" validate:"url"`
}

type UserRes struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
}

type UserDetailRes struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	DisplayName    string `json:"display_name"`
	Bio            string `json:"bio"`
	AvatarURL      string `json:"avatar_url"`
	Role           string `json:"role"`
	FollowerCount  int64  `json:"follower_count"`
	FollowingCount int64  `json:"following_count"`
	PostCount      int64  `json:"post_count"`
	IsFollowing    bool   `json:"is_following"`
	CreatedAt      string `json:"created_at"`
}

type LoginRes struct {
	AccessToken string  `json:"access_token"`
	User        UserRes `json:"user"`
}
