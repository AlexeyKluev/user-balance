package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
)

type AccuralRepo struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func NewAccuralRepo(db *sqlx.DB, logger *zap.Logger) *AccuralRepo {
	return &AccuralRepo{
		logger: logger,
		db:     db,
	}
}

func (r *AccuralRepo) Accural(ctx context.Context, input dto.AccuralDTO) error {
	tx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error("failed to begin database transaction", zap.Error(err))
		return err
	}

	sqlq, args, err := sq.Insert("transactions").
		Columns("user_id", "change").
		Values(input.UserID, input.Amount).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		r.logger.Error("failed to build SQL query", zap.Error(err))
		errTx := tx.Rollback()
		if errTx != nil {
			r.logger.Error("failed to rollback transaction", zap.Error(errTx))
		}
		return err
	}

	_, err = tx.ExecContext(ctx, sqlq, args...)
	if err != nil {
		r.logger.Error("failed to insert transaction", zap.Error(err))
		errTx := tx.Rollback()
		if errTx != nil {
			r.logger.Error("failed to rollback transaction", zap.Error(errTx))
		}

		return err
	}

	args = []interface{}{input.Amount, input.UserID}
	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = users.balance + $1 WHERE id = $2", args...)
	if err != nil {
		r.logger.Error("failed to update user", zap.Error(err))
		errTx := tx.Rollback()
		if errTx != nil {
			r.logger.Error("failed to rollback transaction", zap.Error(errTx))
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		errTx := tx.Rollback()
		if errTx != nil {
			r.logger.Error("failed to rollback transaction", zap.Error(errTx))
		}
		return err
	}

	return nil
}
