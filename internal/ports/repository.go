package ports

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"time"
)

type IMinio interface {
	GenerateUploadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	GenerateDownloadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error
	DownloadFile(ctx context.Context, bucketName, objectName string) (*minio.Object, error)
	GetStatFile(ctx context.Context, bucketName, objectName string) (*minio.ObjectInfo, error)
	DownloadFileRange(ctx context.Context, bucketName, storagePath string, start, end int64) (io.ReadCloser, error)
}

type ICache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiry time.Duration) error
	Delete(ctx context.Context, key string) error
}
