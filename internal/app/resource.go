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
	Config           *config.Config
	Logger           *zap.Logger
	MetricsCollector *metrics.Collector
	UserBalanceUC    *usecase.UserBalanceUseCase
}

func NewResources(config *config.Config) (*Resources, error) {
	// Логгер
	logger, err := newLogger(config.IsProduction)
	if err != nil {
		return nil, err
	}

	// Метрики
	metricsCollector := metrics.NewCollector(version.App)

	// Репозиторий
	repo, err := repository.NewRepository(
		// config.Postgres,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// Services
	userService := service.NewUserService(repo.User)

	// UseCases
	balanceUseCase := usecase.NewUserBalanceUseCase(userService)

	return &Resources{
		Config:           config,
		Logger:           logger,
		MetricsCollector: metricsCollector,
		UserBalanceUC:    balanceUseCase,
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
