package rpc

import (
	mediav1 "github.com/co1seam/ember-backend-api-contracts/gen/go/media"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"github.com/co1seam/ember-backend-media/internal/core/services"
)

type Handler struct {
	Media mediav1.MediaServiceServer
	opts  *models.Options
}

func NewHandler(service *services.Services, opts *models.Options) *Handler {
	return &Handler{
		Media: NewMediaHandler(service.Media, opts),
		opts:  opts,
	}
}
