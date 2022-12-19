package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func (p *PgDB) Master() *sqlx.DB {
	return p.master
}

func (p *PgDB) MasterPing(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := p.master.PingContext(ctx)
	if err != nil {
		p.logger.Error("ping master error: %v", zap.Error(err))
	}

	return err == nil
}

func (p *PgDB) Slave() *sqlx.DB {
	return p.slave
}

func (p *PgDB) SlavePing(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := p.slave.PingContext(ctx)
	if err != nil {
		p.logger.Error("ping slave error: %v", zap.Error(err))
	}
	return err == nil
}

func (p *PgDB) Close() error {
	errMaster, errSlave := p.master.Close(), p.slave.Close()
	if errMaster != nil {
		return errMaster
	}

	return errSlave
}
