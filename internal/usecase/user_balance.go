package usecase

type UserBalanceUseCase struct {
	userService UserService
}

func NewUserBalanceUseCase(us UserService) *UserBalanceUseCase {
	return &UserBalanceUseCase{us}
}

type UserService interface {
	Balance(id int64) (string, error)
}

func (u *UserBalanceUseCase) Balance(id int64) (string, error) {
	balance, err := u.userService.Balance(id)
	if err != nil {
		return "", err
	}

	return balance, nil
}
