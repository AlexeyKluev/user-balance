package postgres

import (
	"github.com/Rhymond/go-money"
	"go.uber.org/zap"

	"github.com/AlexeyKluev/user-balance/internal/domain/model"
)

type UserRepo struct {
	logger *zap.Logger
}

func NewUserRepo(l *zap.Logger) *UserRepo {
	return &UserRepo{
		logger: l,
	}
}

func (r *UserRepo) GetByID(_ int64) (model.User, error) {
	return model.User{
		ID:        213,
		FirstName: "Алексей",
		LastName:  "Клюев",
		Balance:   money.New(100000, money.RUB),
	}, nil
}
