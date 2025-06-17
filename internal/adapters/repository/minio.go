package repository

import (
	"context"
	"github.com/co1seam/ember-backend-media/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &Minio{
		Client: client,
		Bucket: cfg.Bucket,
	}, nil
}

func (m *Minio) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	_, err := m.Client.PutObject(
		ctx,
		bucketName,
		objectName,
		reader,
		objectSize,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	return err
}

func (m *Minio) DownloadFile(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	media, err := m.Client.GetObject(
		ctx,
		bucketName,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (m *Minio) DownloadFileRange(ctx context.Context, bucketName, storagePath string, start, end int64) (io.ReadCloser, error) {
	opts := minio.GetObjectOptions{}
	if start > 0 || end > 0 {
		opts.SetRange(start, end)
	}

	obj, err := m.Client.GetObject(
		ctx,
		bucketName,
		storagePath,
		opts,
	)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (m *Minio) GetStatFile(ctx context.Context, bucketName, objectName string) (*minio.ObjectInfo, error) {
	fileInfo, err := m.Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &fileInfo, nil
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
