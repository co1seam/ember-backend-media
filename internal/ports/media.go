package ports

import (
	"context"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"time"
)

type (
	IMediaRepo interface {
		Create(ctx context.Context, media *models.Media) error
		GetByID(ctx context.Context, id string) (*models.Media, error)
		Update(ctx context.Context, media *models.Media) error
		Delete(ctx context.Context, id string) error
		ListByOwner(ctx context.Context, ownerID string, limit int) ([]*models.Media, error)
	}

	IMediaService interface {
		CreateMedia(ctx context.Context, req *models.CreateMediaRequest) (*models.Media, error)
		GetMedia(ctx context.Context, id string) (*models.Media, error)
		UpdateMedia(ctx context.Context, req *models.UpdateMediaRequest) (*models.Media, error)
		DeleteMedia(ctx context.Context, id string) error
		ListMedia(ctx context.Context, ownerID string, limit int) ([]*models.Media, error)
		GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	}
)
