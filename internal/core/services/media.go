package services

import (
	"context"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"github.com/co1seam/ember-backend-media/internal/ports"
	"github.com/google/uuid"

	"time"
)

type Media struct {
	repo  ports.IMediaRepo
	minio ports.IMinio
	opts  *models.Options
}

func NewMedia(repo ports.IMediaRepo, minio ports.IMinio, opts *models.Options) *Media {
	return &Media{
		repo:  repo,
		minio: minio,
		opts:  opts,
	}
}

func (m *Media) CreateMedia(ctx context.Context, req *models.CreateMediaRequest) (*models.Media, error) {
	media := &models.Media{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		ContentType: req.ContentType,
		StoragePath: "media/" + uuid.New().String() + getFileExtension(req.ContentType),
		OwnerID:     req.OwnerID,
		CreatedAt:   time.Now(),
	}

	if err := m.repo.Create(ctx, media); err != nil {
		return nil, err
	}

	uploadURL, err := m.minio.GenerateUploadURL(ctx, media.StoragePath, 1*time.Hour)
	if err != nil {
		return nil, err
	}
	media.URL = uploadURL

	return media, nil
}

func (m *Media) GetMedia(ctx context.Context, id string) (*models.Media, error) {
	media, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	downloadURL, err := m.minio.GenerateDownloadURL(ctx, media.StoragePath, 24*time.Hour)
	if err != nil {
		return nil, err
	}
	media.URL = downloadURL

	return media, nil
}

func (m *Media) UpdateMedia(ctx context.Context, req *models.UpdateMediaRequest) (*models.Media, error) {
	media, err := m.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	media.Title = req.Title
	media.Description = req.Description

	if err := m.repo.Update(ctx, media); err != nil {
		return nil, err
	}

	return media, nil
}

func (m *Media) DeleteMedia(ctx context.Context, id string) error {
	return m.repo.Delete(ctx, id)
}

func (m *Media) ListMedia(ctx context.Context, ownerID string, limit int) ([]*models.Media, error) {
	return m.repo.ListByOwner(ctx, ownerID, limit)
}

func (m *Media) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	return m.minio.GenerateDownloadURL(ctx, objectName, expiry)
}

func getFileExtension(contentType string) string {
	switch contentType {
	case "video/mp4":
		return ".mp4"
	case "video/mov":
		return ".mov"
	case "video/avi":
		return ".avi"
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ".bin"
	}
}
