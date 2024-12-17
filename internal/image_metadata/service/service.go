package service

import (
	"context"

	"github.com/danushk97/image-analyzer/internal/constants"
	"github.com/danushk97/image-analyzer/internal/image_metadata/dtos"
	"github.com/danushk97/image-analyzer/internal/image_metadata/model/v1"
	"github.com/danushk97/image-analyzer/internal/image_metadata/repo/sql"

	"github.com/danushk97/image-analyzer/pkg/contextkey"
	"github.com/danushk97/image-analyzer/pkg/errors"
)

// Service is offer base service, this is used by all child services of offers
type Service struct {
	Repo *sql.Repo
}

// NewService returns the instance of Service with all options applied
func NewService(opts ...Option) *Service {
	svc := NewOptions(opts...)
	return svc
}

func (s *Service) CreateImageMetadata(
	ctx context.Context,
	req *dtos.CreateImageMetadataRequest,
) (*model.ImageMetadata, errors.IError) {
	imageMetadata := &model.ImageMetadata{
		Filename: req.FileName,
		UserID:   contextkey.GetUserIDFromCtx(ctx),
		Status:   constants.StatusInitiated,
	}
	err := s.Repo.CreateImageMetadata(ctx, imageMetadata)

	if err != nil {
		return imageMetadata, err
	}

	return imageMetadata, nil
}
