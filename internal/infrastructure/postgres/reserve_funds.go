package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rhymond/go-money"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
	"github.com/AlexeyKluev/user-balance/internal/domain/model"
)

type ReserveFundsRepo struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewReserveFundsRepo(db *sqlx.DB, logger *zap.Logger) *ReserveFundsRepo {
	return &ReserveFundsRepo{
		db:     db,
		logger: logger,
	}
}

func (r *ReserveFundsRepo) Reserve(ctx context.Context, input dto.ReservationDTO) error {
	tx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error("failed to begin transaction")
		return err
	}
	balance, err := r.balance(ctx, tx, input.UserID)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			r.logger.Error("failed to rollback transaction", zap.Error(errTx))
		}

		return err
	}

	ok, err := balance.GreaterThan(money.New(input.Amount, money.RUB))
	if err != nil {
		return err
	}

	if !ok {
		return ErrInsufficientBalance
	}

	sqlq, args, err := sq.Insert("reservations").
		SetMap(map[string]interface{}{
			"user_id":    input.UserID,
			"amount":     input.Amount,
			"service_id": input.ServiceID,
			"order_id":   input.OrderID,
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			r.logger.Error("failed to rollback transaction", zap.Error(errTx))
		}

		return err
	}

	_, err = tx.ExecContext(ctx, sqlq, args...)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			r.logger.Error("failed to rollback transaction", zap.Error(errTx))
		}

		return err
	}

	args = []interface{}{input.Amount, input.UserID}
	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = users.balance - $1 WHERE id = $2", args...)
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
		return err
	}

	return nil
}

func (r *ReserveFundsRepo) balance(ctx context.Context, tx *sqlx.Tx, id int64) (*money.Money, error) {
	var user model.User

	sqlq, args, err := sq.Select(
		"id",
		"status",
		"balance",
		"created_at",
		"modified_at",
	).From("users").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	if err = tx.GetContext(ctx, &user, sqlq, args...); err != nil {
		r.logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}

	return money.New(user.Balance, money.RUB), nil
}
