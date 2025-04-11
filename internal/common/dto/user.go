package dto

import "time"

type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatar_url"`
	Locked    bool      `json:"locked"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
