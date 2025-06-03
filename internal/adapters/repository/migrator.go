package repository

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var files embed.FS

type Migrator struct {
	srcDriver source.Driver
}

func NewMigrator() *Migrator {
	src, err := iofs.New(files, "migrations")
	if err != nil {
		log.Errorf("error: %w", err)
	}

	return &Migrator{srcDriver: src}
}

func (m *Migrator) Up(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("embed_migrations", m.srcDriver, "auth", driver)
	if err != nil {
		return err
	}
	defer func() {
		if _, err := migrator.Close(); err != nil {
			return
		}
	}()

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}

	return nil
}

func (m *Migrator) Down(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	operation, err := migrate.NewWithInstance("embed_migrations", m.srcDriver, "postgres", driver)
	if err != nil {
		return err
	}

	defer func() {
		if _, err := operation.Close(); err != nil {
			return
		}
	}()

	if err := operation.Down(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return err
	}

	return nil
}
