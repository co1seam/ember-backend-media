package repository

import (
	"database/sql"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"github.com/co1seam/ember-backend-media/internal/ports"
)

type Repository struct {
	Media ports.IMediaRepo
	MinIO *Minio
	Cache *Redis
}

func NewRepository(db *sql.DB, minio *Minio, cache *Redis, opts *models.Options) *Repository {
	return &Repository{
		Media: NewMedia(db, opts),
		MinIO: minio,
		Cache: cache,
	}
}
