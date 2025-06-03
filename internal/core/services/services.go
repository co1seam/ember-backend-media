package services

import (
	"github.com/co1seam/ember-backend-media/internal/adapters/repository"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"github.com/co1seam/ember-backend-media/internal/ports"
)

type Services struct {
	Media ports.IMediaService
}

func NewService(repos *repository.Repository, opts *models.Options) *Services {
	return &Services{
		Media: NewMedia(repos.Media, repos.MinIO, opts),
	}
}
