package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type embedFS = embed.FS

type EmbedFSWrapper struct {
	embedFS
}

var constants = map[string]string{
	"CORE_PLUGIN_ID":    "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
	"JWT_KEY_SECRET":    umid.New().String(),
	"JWT_KEY_SIGNATURE": umid.New().String(),
	"NODE_ID":           umid.New().String(),
	"BASE_DOMAIN":       "",
}

var constantRegexp = regexp.MustCompile(`\{\{[A-Z,_]*\}\}`)

type CustomFile struct {
	io.ReadCloser
	stat fs.FileInfo
}

func (f CustomFile) Stat() (fs.FileInfo, error) {
	return f.stat, nil
}

//func (e EmbedFSWrapper) ReadFile(name string) ([]byte, error) {
//	fmt.Println(" EmbedFSWrapper ReadFile:" + name)
//	b, err := e.embedFS.ReadFile(name)
//	if err != nil {
//		return nil, err
//	}
//	s := strings.Replace(string(b), "{{BAD_BEGIN}}", "BEGIN", -1)
//	return ([]byte)(s), nil
//}

func (e EmbedFSWrapper) Open(name string) (fs.File, error) {
	fmt.Println(" EmbedFSWrapper Open:" + name)
	f, err := e.embedFS.Open(name)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	s := string(b)

	for key, value := range constants {
		s = strings.Replace(s, "{{"+key+"}}", value, -1)
	}

	unusedConstants := constantRegexp.FindStringSubmatch(s)
	if len(unusedConstants) > 0 {
		names := ""
		for _, n := range unusedConstants {
			names += n + ", "
		}
		return nil, errors.New("Some constants were not replaced:" + names)
	}

	r := io.NopCloser(strings.NewReader(s))

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	customFile := CustomFile{
		ReadCloser: r,
		stat:       stat,
	}

	return customFile, nil
}

func NewEmbedFSWrapper(efs embed.FS) EmbedFSWrapper {
	return EmbedFSWrapper{embedFS: efs}
}

//go:embed sql/*
var migrationFS embed.FS

func pgDBMigrationsConnect(ctx context.Context, log *zap.SugaredLogger, cfg *config.Postgres) (*sql.DB, error) {
	minCfg, err := cfg.MinVersion().GenConfig(log.Desugar())
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate min config")
	}

	db, err := sql.Open("pgx", minCfg.ConnString())
	if err != nil {
		return nil, errors.WithMessage(err, "failed to connect to database")
	}

	if err := db.Ping(); err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if !ok || (ok && pgErr.Code != pgerrcode.InvalidCatalogName) {
			return nil, errors.WithMessage(err, "unknown error with database connection")
		}

		log.Info("Migration: database does not exist")
		if err := createNewDatabase(ctx, log, cfg); err != nil {
			return nil, errors.WithMessage(err, "failed to create new database")
		}

		return pgDBMigrationsConnect(ctx, log, cfg)
	}

	db.SetMaxOpenConns(int(cfg.MAXCONNS))

	return db, nil
}

func createNewDatabase(ctx context.Context, log *zap.SugaredLogger, cfg *config.Postgres) error {
	pgxCfg, err := cfg.GenConfig(log.Desugar())
	if err != nil {
		return errors.WithMessage(err, "failed to generate config")
	}

	cfgNoBase := pgxCfg.ConnConfig
	cfgNoBase.Database = ""

	conn, err := pgx.ConnectConfig(ctx, cfgNoBase)
	if err != nil {
		return errors.WithMessage(err, "unable to connect to database")
	}

	log.Info("Migration: creating database...")
	if _, err := conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %q`, cfg.DATABASE)); err != nil {
		return errors.WithMessage(err, "failed to create database")
	}
	conn.Close(ctx)

	return nil
}

func MigrateDatabase(ctx types.NodeContext, cfg *config.Postgres) error {
	log := ctx.Logger()

	db, err := pgDBMigrationsConnect(ctx, log, cfg)
	if err != nil {
		return errors.WithMessage(err, "failed to create migrations connect")
	}
	defer db.Close()

	config := ctx.Config()
	baseUrl := config.Settings.FrontendURL
	constants["BASE_DOMAIN"] = baseUrl

	wrappedMigrationFS := NewEmbedFSWrapper(migrationFS)
	// get instance of migration data
	data, err := iofs.New(wrappedMigrationFS, "sql")
	if err != nil {
		return errors.WithMessage(err, "failed to get migration data")
	}

	// get DB instance
	pg, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.WithMessage(err, "failed to get database instance")
	}

	// create migration instance
	m, err := migrate.NewWithInstance("iofs", data, "pgx", pg)
	if err != nil {
		return errors.WithMessage(err, "failed to open migration instance")
	}

	var iver uint
	for iver1, _ := data.First(); err == nil; iver1, err = data.Next(iver) {
		iver = iver1
	}

	version, isDirty, err := m.Version()
	log.Infof("Migration: version: %d, %d, %t, %+v", version, iver, isDirty, err)
	if err != nil {
		if err != migrate.ErrNilVersion {
			return errors.WithMessage(err, "failed to obtain current migration version")
		}
		log.Info("Migration: empty (newly created) database detected, will seed!")
	} else if version < iver {
		if isDirty {
			return errors.New("database is dirty")
		}
		log.Infof("Migration: current DB schema verion=%d, available schema version=%d, will migrate", iver, version)
	} else {
		log.Infoln("Migration: migration is not required")
		return nil
	}

	if isDirty {
		return errors.WithMessage(err, "database is dirty, avoiding migration")
	}

	// run your migrations and handle the errors above of course
	if err := m.Up(); err != nil {
		return errors.WithMessage(err, "failed to migrate database")
	}
	version, _, _ = m.Version()
	log.Infof("Migration: success, current schema version: %d", version)

	return nil
}
