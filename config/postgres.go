package config

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v4"
	zapadaptor "github.com/jackc/pgx/v4/log/zapadapter"
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
	x.USERNAME = "postgres"
	x.PORT = 5432
	x.MAXCONNS = 100
}

func (x *Postgres) GenConfig(log *zap.Logger) (*pgxpool.Config, error) {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		x.USERNAME, x.PASSWORD, x.HOST, x.PORT, x.DATABASE)
	if x.MAXCONNS > 0 {
		connString = fmt.Sprintf("%s?pool_max_conns=%d", connString, x.MAXCONNS)
	}

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse postgres config")
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "<unknown>"
	}
	cfg.ConnConfig.RuntimeParams["application_name"] = fmt.Sprintf("Controller on %s", hostname)
	cfg.ConnConfig.LogLevel = pgx.LogLevelError
	cfg.ConnConfig.Logger = zapadaptor.NewLogger(log)

	return cfg, nil
}

func (x *Postgres) MinVersion() *Postgres {
	return &Postgres{
		DATABASE: x.DATABASE,
		HOST:     x.HOST,
		PORT:     x.PORT,
		USERNAME: x.USERNAME,
		PASSWORD: x.PASSWORD,
	}
}
