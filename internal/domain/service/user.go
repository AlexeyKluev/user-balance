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

func (s *UserService) Balance(ctx context.Context, id int64) (string, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()

	user, err := s.repo.GetByID(timeout, id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	return money.New(user.Balance, money.RUB).Display(), nil
}
