package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PgDB struct {
	logger *zap.Logger
	master *sqlx.DB
	slave  *sqlx.DB
}

type DbConnection interface {
	// Master returns master db connection
	Master() *sqlx.DB
	// MasterPing checks master db connection
	MasterPing(ctx context.Context) bool
	// Slave returns slave db connection
	Slave() *sqlx.DB
	// SlavePing checks slave db connection
	SlavePing(ctx context.Context) bool
	// Close close connection
	Close() error
}
