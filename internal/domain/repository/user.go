package repository

import (
	"context"

	"github.com/AlexeyKluev/user-balance/internal/domain/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (model.User, error)
	Create(ctx context.Context, id int64) error
	UserExist(ctx context.Context, id int64) (bool, error)
	UserIsBan(ctx context.Context, id int64) (bool, error)
}
