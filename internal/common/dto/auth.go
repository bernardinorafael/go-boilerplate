package dto

import "time"

type AuthResponse struct {
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserResponse struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
