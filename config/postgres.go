package config

import (
	"fmt"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Postgres struct {
	DATABASE string `yaml:"database" envconfig:"DB_DATABASE"`
	HOST     string `yaml:"host" envconfig:"PGDB_HOST"`
	PORT     uint   `yaml:"port" envconfig:"DB_PORT"`
	USERNAME string `yaml:"username" envconfig:"DB_USERNAME"`
	PASSWORD string `yaml:"password" envconfig:"DB_PASSWORD"`
	MAXCONNS uint   `yaml:"max_conns" envconfig:"DB_MAX_CONNS"`
}

func (x *Postgres) Init() {
	x.DATABASE = "momentum4"
	x.HOST = "localhost"
	x.PASSWORD = ""
	x.USERNAME = "root"
	x.PORT = 5432
	x.MAXCONNS = 100
}

func (x *Postgres) GenConfig() (*pgxpool.Config, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
		x.USERNAME, x.PASSWORD, x.HOST, x.PORT, x.DATABASE, x.MAXCONNS,
	)

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse postgres config")
	}

	return cfg, nil
}

func (x *Postgres) GenMigrateConfig() (*pgx.ConnConfig, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		x.USERNAME, x.PASSWORD, x.HOST, x.PORT, x.DATABASE,
	)

	cfg, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse postgres config")
	}

	return cfg, nil
}
