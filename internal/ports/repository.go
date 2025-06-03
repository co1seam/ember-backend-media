package ports

import (
	"context"
	"time"
)

type IMinio interface {
	GenerateUploadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	GenerateDownloadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
}

type ICache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiry time.Duration) error
	Delete(ctx context.Context, key string) error
}
