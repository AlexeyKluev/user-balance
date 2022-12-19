package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// NewConnection initializes db connection
func NewConnection(config *Config) (*PgDB, error) {
	master, err := initDb(config.Master, config.MasterMaxConns, config.MasterMaxIdleConns, config.MasterConnTTL)
	if err != nil {
		return nil, err
	}

	slave, err := initDb(config.Slave, config.SlaveMaxConns, config.SlaveMaxIdleConns, config.SlaveConnTTL)
	if err != nil {
		return nil, err
	}

	return &PgDB{
		master: master,
		slave:  slave,
	}, nil
}

func initDb(dataSource string, maxConns, maxIdleConns int, connTTL time.Duration) (db *sqlx.DB, err error) {
	db, err = sqlx.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connTTL)

	return
}
