package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/co1seam/ember-backend-media/config"
)

type Postgres struct {
	DB       *sql.DB
	migrator *Migrator
}

func NewPostgres(ctx context.Context, cfg *config.Database) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Name,
		cfg.Pass,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	migrationDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}
	defer migrationDB.Close()

	migrator := NewMigrator()

	if err := migrator.Up(migrationDB); err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	return &Postgres{
		DB:       db,
		migrator: migrator,
	}, nil
}

func (pg *Postgres) Close() error {
	return pg.DB.Close()
}
