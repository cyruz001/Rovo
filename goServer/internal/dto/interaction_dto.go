package dto

type LikePostReq struct {
	PostID string `json:"post_id" validate:"required,uuid4"`
}

type RepostReq struct {
	PostID string `json:"post_id" validate:"required,uuid4"`
}

type FollowUserReq struct {
	FolloweeID string `json:"followee_id" validate:"required,uuid4"`
}

type ErrorRes struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessRes struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
