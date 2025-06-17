package services

import (
	"context"
	"fmt"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"github.com/co1seam/ember-backend-media/internal/ports"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"io"
	"mime"
	"path/filepath"
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
		StoragePath: "",
		OwnerID:     req.OwnerID,
		CreatedAt:   time.Now(),
	}

	if err := m.repo.Create(ctx, media); err != nil {
		return nil, err
	}

	/*
		uploadURL, err := m.minio.GenerateUploadURL(ctx, media.StoragePath, 1*time.Hour)
		if err != nil {
			return nil, err
		}
		media.URL = uploadURL
	*/

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

func (m *Media) UploadFile(ctx context.Context, fileID string, fileName string, size int64, stream io.Reader) (string, error) {
	var media *models.Media
	var err error

	if fileID == "" {
		return "", err
	}

	media, err = m.repo.GetByID(ctx, fileID)
	if err != nil {
		return "", err
	}

	objectPath := fmt.Sprintf("%s/%s/%s", media.OwnerID, media.ID, fileName)

	media.StoragePath = objectPath
	media.URL = fmt.Sprintf("%s%s%s", m.opts.Config.MinIO.Endpoint, m.opts.Config.MinIO.Bucket, objectPath)

	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	media.ContentType = contentType

	if err := m.minio.UploadFile(ctx, m.opts.Config.MinIO.Bucket, objectPath, stream, size, contentType); err != nil {
		return "", err
	}

	if err := m.repo.Update(ctx, media); err != nil {
		return "", err
	}

	return objectPath, nil
}

func (m *Media) DownloadFile(ctx context.Context, fileID string) (*minio.Object, error) {
	return m.minio.DownloadFile(ctx, m.opts.Config.MinIO.Bucket, fileID)
}

func (m *Media) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	return m.minio.GenerateDownloadURL(ctx, objectName, expiry)
}

func (m *Media) GetStatFile(ctx context.Context, objectName string) (*minio.ObjectInfo, error) {
	return m.minio.GetStatFile(ctx, m.opts.Config.MinIO.Bucket, objectName)
}

func (m *Media) DownloadFileRange(ctx context.Context, objectName string, start, end int64) (io.ReadCloser, error) {
	return m.minio.DownloadFileRange(ctx, m.opts.Config.MinIO.Bucket, objectName, start, end)
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
