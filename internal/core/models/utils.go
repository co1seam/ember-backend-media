package models

import (
	"github.com/co1seam/ember-backend-media/config"
	"log/slog"
)

type Options struct {
	Logger *slog.Logger
	Config *config.Config
}

const MaxFileSize = 100 << 20

const (
	MediaTable = "media"
)
