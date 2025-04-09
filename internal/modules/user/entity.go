package user

import "time"

type Entity struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatar_url"`
	Enabled   bool      `json:"enabled"`
	Locked    bool      `json:"locked"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
