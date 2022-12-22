package model

import (
	"time"
)

type User struct {
	ID         int64     `db:"id"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	Status     string    `db:"status"`
	Balance    int64     `db:"balance"`
	CreatedAt  time.Time `db:"created_at"`
	ModifiedAt time.Time `db:"modified_at"`
}
