package usecase

import (
	"context"
	"errors"

	"github.com/AlexeyKluev/user-balance/internal/domain/service"
)

type UserBalanceUseCase struct {
	userService UserService
}

func NewUserBalanceUseCase(us UserService) *UserBalanceUseCase {
	return &UserBalanceUseCase{us}
}

type UserService interface {
	Balance(ctx context.Context, id int64) (string, error)
}

func (u *UserBalanceUseCase) Balance(ctx context.Context, id int64) (string, error) {
	balance, err := u.userService.Balance(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	return balance, nil
}
