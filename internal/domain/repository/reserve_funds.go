package repository

import (
	"context"

	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
)

type ReserveFundsRepository interface {
	Reserve(ctx context.Context, input dto.ReservationDTO) error
}
