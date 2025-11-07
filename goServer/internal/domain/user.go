package domain

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt int64  `json:"created_at"`
}
