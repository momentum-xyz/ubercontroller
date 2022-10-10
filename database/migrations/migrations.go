package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/pkg/errors"
	"log"
	"os"
	"regexp"
)

//go:embed sql/*
var migrationFS embed.FS

func pgDBMigrationsConnect(cfg *pgx.ConnConfig) (*sql.DB, error) {
	cs := cfg.ConnString()
	reg := regexp.MustCompile(`pool_max_conns=[0-9]*`)
	cs = reg.ReplaceAllString(cs, "")
	db, err := sql.Open("pgx", cs)

	if err != nil {
		return nil, errors.WithMessage(err, "Unable to connect to database server")
	}
	err = db.Ping()

	if err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if !ok || (ok && pgErr.Code != pgerrcode.InvalidCatalogName) {
			return nil, errors.WithMessage(err, "Unknown error wth database connection")
		}
		fmt.Println("Database does not exist")
		err = createNewDatabase(cfg)
		if err != nil {
			return nil, err
		}
		return pgDBMigrationsConnect(cfg)
	}
	db.SetMaxOpenConns(100)
	return db, nil
}

func createNewDatabase(cfg *pgx.ConnConfig) error {
	cfgNoBase := cfg.Copy()
	cfgNoBase.Database = ""

	conn, err := pgx.ConnectConfig(context.Background(), cfgNoBase)
	if err != nil {
		return errors.WithMessage(err, "Unable to connect to database server")
	}

	fmt.Println("Creating database...")
	_, err = conn.Exec(context.Background(), "CREATE DATABASE "+cfg.Database)
	if err != nil {
		return errors.WithMessage(err, "Can not create database")
	}
	conn.Close(context.Background())
	return nil
}

func MigrateDatabase(cfg *config.Postgres) error {
	config, err := cfg.GenMigrateConfig()
	db, err := pgDBMigrationsConnect(config)
	//fmt.Println("MIGRATE DATABASE")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	// get instance of migration data

	d, err := iofs.New(migrationFS, "sql")
	if err != nil {
		log.Fatal(err)
	}
	// get DB instance
	pg, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.WithMessage(err, "Can not get DB instance")
	}

	// create migration instance
	m, err := migrate.NewWithInstance("iofs", d, "pgx", pg)
	if err != nil {
		return errors.WithMessage(err, "Can not open migration instance")
	}

	var iver uint
	for iver1, _ := d.First(); err == nil; iver1, err = d.Next(iver) {
		iver = iver1
	}

	version, isDirty, err := m.Version()
	fmt.Println("Version:", version, iver, isDirty, err)
	if err != nil {
		if err != migrate.ErrNilVersion {
			return errors.WithMessage(err, "Can not obtain current migration version")
		}
		fmt.Println("Empty (newly created) database detected, will seed!")
	} else if version < iver {
		if isDirty {
			return errors.New("Database is dirty")
		}
		fmt.Printf("Current DB schema verion=%d, available schema version=%d, will miigrate\n", iver, version)
	} else {
		fmt.Println("No DB migration required")
		return nil
	}

	if isDirty {
		return errors.WithMessage(err, "Database is dirty, avoiding migration")
	}
	//fmt.Println("Version:", version, isDirty, err)

	err = m.Up() // run your migrations and handle the errors above of course
	fmt.Println("ff", err)
	if err != nil {
		return errors.WithMessage(err, "Migration failed")
	}
	fmt.Println("ff1")
	version, _, _ = m.Version()
	fmt.Println("ff2")
	fmt.Println("Migration was successful, current schema version:", version)

	return nil
}
