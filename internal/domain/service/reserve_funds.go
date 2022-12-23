package service

import (
	"context"

	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
	"github.com/AlexeyKluev/user-balance/internal/domain/repository"
)

type ReserveFundsService struct {
	rfr repository.ReserveFundsRepository
}

func NewReserveFundsService(rfr repository.ReserveFundsRepository) *ReserveFundsService {
	return &ReserveFundsService{
		rfr: rfr,
	}
}

func (s *ReserveFundsService) ReserveFunds(ctx context.Context, input dto.ReservationDTO) error {
	return s.rfr.Reserve(ctx, input)
}
