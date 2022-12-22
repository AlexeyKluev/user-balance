package usecase

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

var ErrUserIsBanned = errors.New("user is banned")
