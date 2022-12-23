package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"

	"github.com/AlexeyKluev/user-balance/internal/domain/enum"
	"github.com/AlexeyKluev/user-balance/internal/domain/model"
)

type UserRepo struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func NewUserRepo(db *sqlx.DB, l *zap.Logger) *UserRepo {
	return &UserRepo{
		logger: l,
		db:     db,
	}
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (model.User, error) {
	var user model.User

	sqlq, args, err := sq.Select(
		"id",
		"status",
		"balance",
		"created_at",
		"modified_at",
	).From("users").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return model.User{}, err
	}

	err = r.db.GetContext(ctx, &user, sqlq, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, ErrNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepo) UserExist(ctx context.Context, id int64) (bool, error) {
	var user model.User

	sqlq, args, err := sq.Select(
		"id",
		"status",
		"balance",
		"created_at",
		"modified_at",
	).
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return false, err
	}

	err = r.db.GetContext(ctx, &user, sqlq, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *UserRepo) UserIsBan(ctx context.Context, id int64) (bool, error) {
	var user model.User

	sqlq, args, err := sq.Select(
		"id",
		"status",
		"balance",
		"created_at",
		"modified_at",
	).
		From("users").
		Where(sq.Eq{"id": id, "status": enum.StatusActive}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return false, err
	}

	err = r.db.GetContext(ctx, &user, sqlq, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return true, err
	}

	return false, nil
}

func (r *UserRepo) Create(ctx context.Context, id int64) error {
	sqlq, args, err := sq.Insert("users").
		Columns("id", "status", "balance").
		Values(id, enum.StatusActive, 0).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, sqlq, args...); err != nil {
		return err
	}

	return nil
}
