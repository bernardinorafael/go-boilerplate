package dto

import "time"

type AuthResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `jsono:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
