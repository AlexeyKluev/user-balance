package usecase

import (
	"context"

	"github.com/Rhymond/go-money"

	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
)

type AccuralFundsService interface {
	Accural(context.Context, dto.AccuralDTO) error
}

type UserBanService interface {
	UserIsBan(ctx context.Context, id int64) (bool, error)
}

type UserExistService interface {
	UserExist(ctx context.Context, id int64) (bool, error)
}

type UserBalanceService interface {
	Balance(ctx context.Context, id int64) (*money.Money, error)
}

type ReserveFundsService interface {
	ReserveFunds(ctx context.Context, input dto.ReservationDTO) error
}

type UserCreateService interface {
	Create(ctx context.Context, id int64) error
}
