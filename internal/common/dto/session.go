package dto

import "time"

type LoginResponse struct {
	SessionID    string `json:"session_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateSession struct {
	UserID       string `json:"user_id"`
	IP           string `json:"ip"`
	Agent        string `json:"agent"`
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessToken struct {
	AccessToken        string    `json:"access_token"`
	AccessTokenExpires time.Time `json:"access_token_expires"`
}

type SessionResponse struct {
	ID      string    `json:"id"`
	Agent   string    `json:"agent"`
	IP      string    `json:"ip_address"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
