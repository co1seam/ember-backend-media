package rpc

import (
	"context"
	mediav1 "github.com/co1seam/ember-backend-api-contracts/gen/go/media"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"github.com/co1seam/ember-backend-media/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"os"
	"path/filepath"
	"time"
)

type MediaHandler struct {
	mediav1.UnimplementedMediaServiceServer
	service ports.IMediaService
	opts    *models.Options
}

type chanReader struct {
	data <-chan []byte
	err  <-chan error
	sem  chan<- struct{}
	ctx  context.Context
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
		Media: toProtoMedia(media, 0),
	}, nil
}

func (h *MediaHandler) GetMedia(ctx context.Context, req *mediav1.GetMediaRequest) (*mediav1.MediaResponse, error) {
	media, err := h.service.GetMedia(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "media not found")
	}

	mediaInfo, err := h.service.GetStatFile(ctx, media.StoragePath)
	if err != nil {
		return nil, status.Error(codes.NotFound, "media not found")
	}

	return &mediav1.MediaResponse{
		Media: toProtoMedia(media, mediaInfo.Size),
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
		Media: toProtoMedia(media, -1),
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
		mediaInfo, err := h.service.GetStatFile(ctx, media.ID)
		if err != nil {
			return nil, status.Error(codes.Internal, "list failed")
		}
		protoMedia[i] = toProtoMedia(media, mediaInfo.Size)
	}

	return &mediav1.ListMediaResponse{
		Media: protoMedia,
	}, nil
}

func toProtoMedia(media *models.Media, size int64) *mediav1.Media {
	return &mediav1.Media{
		Id:          media.ID,
		Title:       media.Title,
		Description: media.Description,
		ContentType: media.ContentType,
		OwnerId:     media.OwnerID,
		CreatedAt:   media.CreatedAt.Format(time.RFC3339),
		Url:         media.URL,
		Size:        size,
	}
}

func (h *MediaHandler) UploadFile(stream mediav1.MediaService_UploadFileServer) error {
	var FileID string
	var fileName string
	var tempFile *os.File
	var totalSize int64
	var err error

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if chunk.IsFirst {
			FileID = chunk.FileId
			fileName = chunk.FileName
			totalSize = chunk.TotalSize
			tempFile, err = os.CreateTemp("", filepath.Base(fileName)+"-*")
			if err != nil {
				return status.Error(codes.Internal, "failed to create temp file")
			}
			defer os.Remove(tempFile.Name())
			defer tempFile.Close()
			continue
		}

		if _, err := tempFile.Write(chunk.Content); err != nil {
			return status.Error(codes.Internal, "write error: "+err.Error())
		}
	}

	if _, err := tempFile.Seek(0, 0); err != nil {
		return status.Error(codes.Internal, "seek error: "+err.Error())
	}

	Url, err := h.service.UploadFile(stream.Context(), FileID, fileName, totalSize, tempFile)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return stream.SendAndClose(&mediav1.FileResponse{
		FileId: FileID,
		Url:    Url,
	})
}

func (h *MediaHandler) DownloadFile(req *mediav1.FileRequest, stream mediav1.MediaService_DownloadFileServer) error {
	meta, err := h.service.GetMedia(stream.Context(), req.FileId)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	fileInfo, err := h.service.GetStatFile(stream.Context(), meta.StoragePath)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if meta.OwnerID != req.OwnerId {
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	start := req.Start
	end := req.End

	if end < 0 || end >= fileInfo.Size {
		end = fileInfo.Size - 1
	}

	contentLength := end - start + 1
	if contentLength <= 0 {
		return status.Error(codes.InvalidArgument, "invalid range")
	}

	reader, err := h.service.DownloadFileRange(stream.Context(), meta.StoragePath, start, end)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer reader.Close()

	buf := make([]byte, 64*1024)
	var totalSent int64

	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return status.Error(codes.Internal, err.Error())
		}

		if n > 0 {
			if err := stream.Send(&mediav1.DownloadResponse{Chunk: buf[:n]}); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
			totalSent += int64(n)
		}

		// Проверяем завершение чтения
		if err == io.EOF || totalSent >= contentLength {
			break
		}
	}

	return nil
}

func (r *chanReader) Read(p []byte) (int, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	case err := <-r.err:
		return 0, err
	case data, ok := <-r.data:
		if !ok {
			return 0, io.EOF
		}

		go func() {
			select {
			case r.sem <- struct{}{}:
			case <-r.ctx.Done():
			}
		}()

		n := copy(p, data)
		return n, nil
	}
}
