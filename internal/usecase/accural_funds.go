package usecase

import (
	"context"
	"time"

	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
)

type AccuralFundsUseCase struct {
	accuralFundsService AccuralFundsService
	userBanService      UserBanService
	userExistService    UserExistService
	userCreateService   UserCreateService
}

func NewAccuralFundsUseCase(
	afs AccuralFundsService,
	ubs UserBanService,
	ues UserExistService,
	ucs UserCreateService,
) *AccuralFundsUseCase {
	return &AccuralFundsUseCase{
		accuralFundsService: afs,
		userBanService:      ubs,
		userExistService:    ues,
		userCreateService:   ucs,
	}
}

func (uc *AccuralFundsUseCase) Accural(ctx context.Context, input dto.AccuralDTO) error {
	timeoutUserExist, cancelFuncUserExist := context.WithTimeout(ctx, time.Second*1)
	defer cancelFuncUserExist()
	exist, err := uc.userExistService.UserExist(timeoutUserExist, input.UserID)
	if err != nil {
		return err
	}

	if !exist {
		if err := uc.userCreateService.Create(ctx, input.UserID); err != nil {
			return err
		}
	}

	timeoutUserIsBan, cancelFuncUserIsBan := context.WithTimeout(ctx, time.Second*1)
	defer cancelFuncUserIsBan()
	isBan, err := uc.userBanService.UserIsBan(timeoutUserIsBan, input.UserID)
	if err != nil {
		return err
	}
	if isBan {
		return ErrUserIsBanned
	}

	timeoutAccural, cancelFuncAccural := context.WithTimeout(ctx, time.Second*3)
	defer cancelFuncAccural()
	err = uc.accuralFundsService.Accural(timeoutAccural, input)
	if err != nil {
		return err
	}

	return nil
}
