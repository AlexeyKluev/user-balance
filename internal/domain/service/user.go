package service

import (
	"context"
	"errors"
	"time"

	"github.com/Rhymond/go-money"

	"github.com/AlexeyKluev/user-balance/internal/domain/repository"
	"github.com/AlexeyKluev/user-balance/internal/infrastructure/postgres"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) Balance(ctx context.Context, id int64) (*money.Money, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()

	user, err := s.repo.GetByID(timeout, id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return money.New(user.Balance, money.RUB), nil
}

func (s *UserService) UserIsBan(ctx context.Context, id int64) (bool, error) {
	isBan, err := s.repo.UserIsBan(ctx, id)
	if err != nil {
		return true, err
	}

	return isBan, nil
}

func (s *UserService) UserExist(ctx context.Context, id int64) (bool, error) {
	isExist, err := s.repo.UserExist(ctx, id)
	if err != nil {
		return false, err
	}

	return isExist, nil
}
