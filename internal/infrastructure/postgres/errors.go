package postgres

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

var ErrInsufficientBalance = errors.New("insufficient balance")
