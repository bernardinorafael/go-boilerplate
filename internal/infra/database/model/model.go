package model

import "time"

type User struct {
	ID      string    `db:"id"`
	Name    string    `db:"name"`
	Email   string    `db:"email"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

type Product struct {
	ID      string    `db:"id"`
	Name    string    `db:"name"`
	Price   int64     `db:"price"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

type Code struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Code      string     `db:"code"`
	Active    bool       `db:"active"`
	Attempts  int64      `db:"attempts"`
	UsedAt    *time.Time `db:"used_at"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}
