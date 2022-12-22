package usecase

import (
	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
)

type AccuralFundsUseCase struct {
	accuralFundsService AccuralFundsService
	userBanService      UserBanService
	userExistService    UserExistService
}

func NewAccuralFundsUseCase(
	afs AccuralFundsService,
	ubs UserBanService,
	ues UserExistService,
) *AccuralFundsUseCase {
	return &AccuralFundsUseCase{
		accuralFundsService: afs,
		userBanService:      ubs,
		userExistService:    ues,
	}
}

type AccuralFundsService interface {
	Accural(dto.AccuralDTO) error
}

type UserBanService interface {
	UserIsBan(id int64) (bool, error)
}

type UserExistService interface {
	UserExist(id int64) (bool, error)
}

func (uc *AccuralFundsUseCase) Accural(input dto.AccuralDTO) error {
	exist, err := uc.userExistService.UserExist(input.UserID)
	if err != nil {
		return err
	}
	if !exist {
		return ErrNotFound
	}

	isBan, err := uc.userBanService.UserIsBan(input.UserID)
	if err != nil {
		return err
	}
	if isBan {
		return ErrUserIsBanned
	}

	err = uc.accuralFundsService.Accural(input)
	if err != nil {
		return err
	}

	return nil
}
