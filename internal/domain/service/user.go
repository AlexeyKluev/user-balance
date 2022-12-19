package service

import (
	"github.com/AlexeyKluev/user-balance/internal/domain/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) Balance(id int64) (string, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}

	return user.Balance.Display(), nil
}
