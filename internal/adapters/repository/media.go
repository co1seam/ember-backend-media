package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"github.com/co1seam/ember-backend-media/internal/ports"
)

type Media struct {
	db   *sql.DB
	opts *models.Options
}

func NewMedia(db *sql.DB, opts *models.Options) ports.IMediaRepo {
	return &Media{
		db:   db,
		opts: opts,
	}
}

func (m *Media) Create(ctx context.Context, media *models.Media) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (id, title, description, content_type, storage_path, owner_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		models.MediaTable,
	)

	_, err := m.db.ExecContext(
		ctx,
		query,
		media.ID,
		media.Title,
		media.Description,
		media.ContentType,
		media.StoragePath,
		media.OwnerID,
		media.CreatedAt,
	)
	return err
}

func (m *Media) GetByID(ctx context.Context, id string) (*models.Media, error) {
	query := fmt.Sprintf(
		"SELECT id, title, description, content_type, storage_path, owner_id, created_at FROM %s WHERE id = $1",
		models.MediaTable,
	)

	row := m.db.QueryRowContext(ctx, query, id)

	media := &models.Media{}
	err := row.Scan(
		&media.ID,
		&media.Title,
		&media.Description,
		&media.ContentType,
		&media.StoragePath,
		&media.OwnerID,
		&media.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return media, nil
}

func (m *Media) Update(ctx context.Context, media *models.Media) error {
	query := fmt.Sprintf(
		"UPDATE %s SET title = $1, description = $2, content_type = $3, storage_path = $4, owner_id = $5, created_at = $6 WHERE id = $7",
		models.MediaTable,
	)

	_, err := m.db.ExecContext(
		ctx,
		query,
		media.Title,
		media.Description,
		media.ContentType,
		media.StoragePath,
		media.OwnerID,
		media.CreatedAt,
		media.ID,
	)
	return err
}

func (m *Media) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id = $1",
		models.MediaTable,
	)

	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m *Media) ListByOwner(ctx context.Context, ownerID string, limit int) ([]*models.Media, error) {
	query := fmt.Sprintf(
		"SELECT id, title, description, content_type, storage_path, owner_id, created_at FROM %s WHERE owner_id = $1 ORDER BY created_at DESC LIMIT $2",
		models.MediaTable,
	)

	rows, err := m.db.QueryContext(ctx, query, ownerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediaList []*models.Media
	for rows.Next() {
		media := &models.Media{}
		err := rows.Scan(
			&media.ID,
			&media.Title,
			&media.Description,
			&media.ContentType,
			&media.StoragePath,
			&media.OwnerID,
			&media.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		mediaList = append(mediaList, media)
	}

	return mediaList, nil
}
