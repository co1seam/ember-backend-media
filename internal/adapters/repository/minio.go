package repository

import (
	"context"
	"github.com/co1seam/ember-backend-media/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

type Minio struct {
	Client *minio.Client
	Bucket string
}

func NewMinio(cfg *config.MinIO) (*Minio, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Minio{
		Client: client,
		Bucket: cfg.Bucket,
	}, nil
}

func (m *Minio) GenerateUploadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := m.Client.PresignedPutObject(ctx, m.Bucket, objectName, expiry)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (m *Minio) GenerateDownloadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := m.Client.PresignedGetObject(ctx, m.Bucket, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
