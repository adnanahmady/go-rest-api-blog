package database

import (
	"errors"
	"fmt"

	"github.com/adnanahmady/go-rest-api-blog/config"
	"github.com/adnanahmady/go-rest-api-blog/pkg/applog"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseManager interface {
	GetClient() *sqlx.DB
}

type DatabaseManagerImpl struct {
	cfg *config.Config
	db  *sqlx.DB
	lgr applog.Logger
}

func NewDatabaseManager(cfg *config.Config, lgr applog.Logger) *DatabaseManagerImpl {
	return &DatabaseManagerImpl{
		cfg: cfg,
		db:  connect(cfg, lgr),
		lgr: lgr,
	}
}

func connect(cfg *config.Config, lgr applog.Logger) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", cfg.Database.Path)
	if err != nil {
		lgr.Fatal("failed to connect to database", err)
		return nil
	}
	lgr.With("path", cfg.Database.Path).Info("connected to database")
	return db
}

func (m *DatabaseManagerImpl) Close() error {
	return m.db.Close()
}

func (m *DatabaseManagerImpl) Migrate() error {
	mgr, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("sqlite3://%s", m.cfg.Database.Path),
	)

	if err != nil {
		m.lgr.Fatal("failed to create migration manager", err)
		return err
	}

	if err := mgr.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		m.lgr.Fatal("failed to run migrations", err)
		return err
	}

	m.lgr.Info("database migrations applied successfully")
	return nil
}

func (m *DatabaseManagerImpl) GetClient() *sqlx.DB {
	return m.db
}
