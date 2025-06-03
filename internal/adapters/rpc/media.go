package rpc

import (
	"context"
	mediav1 "github.com/co1seam/ember-backend-api-contracts/gen/go/media"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"github.com/co1seam/ember-backend-media/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type MediaHandler struct {
	mediav1.UnimplementedMediaServiceServer
	service ports.IMediaService
	opts    *models.Options
}

func NewMediaHandler(service ports.IMediaService, opts *models.Options) *MediaHandler {
	return &MediaHandler{
		service: service,
		opts:    opts,
	}
}

func (h *MediaHandler) CreateMedia(ctx context.Context, req *mediav1.CreateMediaRequest) (*mediav1.MediaResponse, error) {
	media, err := h.service.CreateMedia(ctx, &models.CreateMediaRequest{
		Title:       req.Title,
		Description: req.Description,
		ContentType: req.ContentType,
		OwnerID:     req.OwnerId,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &mediav1.MediaResponse{
		Media: toProtoMedia(media),
	}, nil
}

func (h *MediaHandler) GetMedia(ctx context.Context, req *mediav1.GetMediaRequest) (*mediav1.MediaResponse, error) {
	media, err := h.service.GetMedia(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "media not found")
	}

	return &mediav1.MediaResponse{
		Media: toProtoMedia(media),
	}, nil
}

func (h *MediaHandler) UpdateMedia(ctx context.Context, req *mediav1.UpdateMediaRequest) (*mediav1.MediaResponse, error) {
	media, err := h.service.UpdateMedia(ctx, &models.UpdateMediaRequest{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "update failed")
	}

	return &mediav1.MediaResponse{
		Media: toProtoMedia(media),
	}, nil
}

func (h *MediaHandler) DeleteMedia(ctx context.Context, req *mediav1.DeleteMediaRequest) (*emptypb.Empty, error) {
	if err := h.service.DeleteMedia(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, "delete failed")
	}
	return &emptypb.Empty{}, nil
}

func (h *MediaHandler) ListMedia(ctx context.Context, req *mediav1.ListMediaRequest) (*mediav1.ListMediaResponse, error) {
	mediaList, err := h.service.ListMedia(ctx, req.OwnerId, int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, "list failed")
	}

	protoMedia := make([]*mediav1.Media, len(mediaList))
	for i, media := range mediaList {
		protoMedia[i] = toProtoMedia(media)
	}

	return &mediav1.ListMediaResponse{
		Media: protoMedia,
	}, nil
}

func toProtoMedia(media *models.Media) *mediav1.Media {
	return &mediav1.Media{
		Id:          media.ID,
		Title:       media.Title,
		Description: media.Description,
		ContentType: media.ContentType,
		OwnerId:     media.OwnerID,
		CreatedAt:   media.CreatedAt.Format(time.RFC3339),
		Url:         media.URL,
	}
}
