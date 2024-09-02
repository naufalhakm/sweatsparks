package params

import "time"

type UserRegisterResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type GetAllUser struct {
	ID        uint64    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
