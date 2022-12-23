package repository

import (
	"context"

	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
)

type AccuralFundsRepository interface {
	Accural(context.Context, dto.AccuralDTO) error
}
