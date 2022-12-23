package service

import (
	"context"

	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
	"github.com/AlexeyKluev/user-balance/internal/domain/repository"
)

type AccuralService struct {
	afr repository.AccuralFundsRepository
}

func NewAccrualService(afr repository.AccuralFundsRepository) *AccuralService {
	return &AccuralService{
		afr: afr,
	}
}

func (s *AccuralService) Accural(ctx context.Context, input dto.AccuralDTO) error {
	return s.afr.Accural(ctx, input)
}
