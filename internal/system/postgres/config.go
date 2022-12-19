package postgres

import (
	"time"
)

// Config postgres configuration
type Config struct {
	Master             string        `required:"true"`
	Slave              string        `required:"true"`
	MasterConnTTL      time.Duration `envconfig:"default=59s"`
	MasterMaxConns     int           `envconfig:"default=0"`
	MasterMaxIdleConns int           `envconfig:"default=2"`
	SlaveConnTTL       time.Duration `envconfig:"default=59s"`
	SlaveMaxConns      int           `envconfig:"default=0"`
	SlaveMaxIdleConns  int           `envconfig:"default=2"`
}
