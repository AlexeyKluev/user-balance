package app

import (
	"go.uber.org/zap"

	"github.com/AlexeyKluev/user-balance/internal/domain/repository"
	"github.com/AlexeyKluev/user-balance/internal/domain/service"
	"github.com/AlexeyKluev/user-balance/internal/metrics"
	"github.com/AlexeyKluev/user-balance/internal/usecase"

	"github.com/AlexeyKluev/user-balance/internal/config"
	"github.com/AlexeyKluev/user-balance/internal/version"
)

type Resources struct {
	Config             *config.Config
	Logger             *zap.Logger
	MetricsCollector   *metrics.Collector
	UserBalanceUC      *usecase.UserBalanceUseCase
	AccuralFundsUC     *usecase.AccuralFundsUseCase
	ReservationFundsUC *usecase.ReservationFundsUseCase
}

func NewResources(config *config.Config) (*Resources, error) {
	// Logger
	logger, err := newLogger(config.IsProduction)
	if err != nil {
		return nil, err
	}

	// Metrics
	metricsCollector := metrics.NewCollector(version.App)

	// Repository
	repo, err := repository.NewRepository(
		config.Postgres,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// Services
	userService := service.NewUserService(repo.User)
	accrualService := service.NewAccrualService(repo.Accural)
	reserveFundsService := service.NewReserveFundsService(repo.Reserve)

	// UseCases
	balanceUseCase := usecase.NewUserBalanceUseCase(userService)
	accuralFundsUseCase := usecase.NewAccuralFundsUseCase(accrualService, userService, userService, userService)
	reservationFundsUseCase := usecase.NewReservationFundsUseCase(userService, userService, reserveFundsService)

	return &Resources{
		Config:             config,
		Logger:             logger,
		MetricsCollector:   metricsCollector,
		UserBalanceUC:      balanceUseCase,
		AccuralFundsUC:     accuralFundsUseCase,
		ReservationFundsUC: reservationFundsUseCase,
	}, nil
}

func (r *Resources) Close() error {
	return nil
}

func newLogger(isProduction bool) (logger *zap.Logger, err error) {
	if isProduction {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	return
}
