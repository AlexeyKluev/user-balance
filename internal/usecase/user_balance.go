package usecase

import (
	"context"
	"errors"

	"github.com/AlexeyKluev/user-balance/internal/domain/service"
)

type UserBalanceUseCase struct {
	userService UserBalanceService
}

func NewUserBalanceUseCase(us UserBalanceService) *UserBalanceUseCase {
	return &UserBalanceUseCase{us}
}

type UserBalanceService interface {
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
