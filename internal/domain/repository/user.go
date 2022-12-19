package repository

import (
	"github.com/AlexeyKluev/user-balance/internal/domain/model"
)

type UserRepository interface {
	GetByID(id int64) (model.User, error)
}
