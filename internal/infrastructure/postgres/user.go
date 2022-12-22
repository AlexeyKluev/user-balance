package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"

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
		"first_name",
		"last_name",
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
