package service

import (
	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
)

type AccuralService struct {
}

func NewAccrualService() *AccuralService {
	return &AccuralService{}
}

func (s *AccuralService) Accural(dto.AccuralDTO) error {
	return nil
}
