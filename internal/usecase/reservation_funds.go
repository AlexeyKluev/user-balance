package usecase

import (
	"context"
	"errors"

	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
	"github.com/AlexeyKluev/user-balance/internal/infrastructure/postgres"
)

type ReservationFundsUseCase struct {
	ues   UserExistService
	ubans UserBanService
	rfs   ReserveFundsService
}

func NewReservationFundsUseCase(
	ues UserExistService,
	ubans UserBanService,
	rfs ReserveFundsService,
) *ReservationFundsUseCase {
	return &ReservationFundsUseCase{
		ues:   ues,
		ubans: ubans,
		rfs:   rfs,
	}
}

func (uc *ReservationFundsUseCase) Reservation(
	ctx context.Context,
	input dto.ReservationDTO,
) error {
	exist, err := uc.ues.UserExist(ctx, input.UserID)
	if err != nil {
		return err
	}

	if !exist {
		return ErrNotFound
	}

	isBan, err := uc.ubans.UserIsBan(ctx, input.UserID)
	if err != nil {
		return err
	}

	if isBan {
		return ErrUserIsBanned
	}

	if err = uc.rfs.ReserveFunds(ctx, input); err != nil {
		if errors.Is(err, postgres.ErrInsufficientBalance) {
			return ErrInsufficientBalance
		}
		return err
	}

	return nil
}
