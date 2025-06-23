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
