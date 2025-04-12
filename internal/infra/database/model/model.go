package model

import "time"

type Session struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	IP           string    `db:"ip_address"`
	Agent        string    `db:"agent"`
	RefreshToken string    `db:"refresh_token"`
	Active       bool      `db:"active"`
	Expires      time.Time `db:"expires"`
	Created      time.Time `db:"created"`
	Updated      time.Time `db:"updated"`
}

type User struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	AvatarURL *string   `db:"avatar_url"`
	Enabled   bool      `db:"enabled"`
	Locked    bool      `db:"locked"`
	Created   time.Time `db:"created"`
	Updated   time.Time `db:"updated"`
}

type Product struct {
	ID      string    `db:"id"`
	Name    string    `db:"name"`
	Price   int64     `db:"price"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}
