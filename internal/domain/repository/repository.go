package repository

import (
	"go.uber.org/zap"

	pgRepo "github.com/AlexeyKluev/user-balance/internal/infrastructure/postgres"
	"github.com/AlexeyKluev/user-balance/internal/system/postgres"
)

type Repository struct {
	pgDB    *postgres.PgDB
	logger  *zap.Logger
	User    UserRepository
	Accural AccuralFundsRepository
	Reserve ReserveFundsRepository
}

func NewRepository(
	config *postgres.Config,
	logger *zap.Logger,
) (*Repository, error) {
	pgDB, err := postgres.NewConnection(config)
	if err != nil {
		return nil, err
	}

	return &Repository{
		pgDB:    pgDB,
		logger:  logger,
		User:    pgRepo.NewUserRepo(pgDB.Master(), logger),
		Accural: pgRepo.NewAccuralRepo(pgDB.Master(), logger),
		Reserve: pgRepo.NewReserveFundsRepo(pgDB.Master(), logger),
	}, nil
}

func (r *Repository) Close() error {
	return r.pgDB.Close()
}
