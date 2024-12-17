package image_metadata

import (
	"context"
	"net/http"

	"github.com/danushk97/image-analyzer/internal/image_metadata/dtos"
	"github.com/danushk97/image-analyzer/internal/image_metadata/service"
	"github.com/danushk97/image-analyzer/internal/middlewares"
	pkgLogger "github.com/danushk97/image-analyzer/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ImageMetadataServer struct {
	service *service.Service
}

// NewServer creates a new server
func NewServer(
	imageMetaService *service.Service,
) *ImageMetadataServer {
	return &ImageMetadataServer{
		service: imageMetaService,
	}
}

func (is *ImageMetadataServer) SetupRoutes(r *gin.Engine) {
	imageApi := r.Group("/v1/images")

	imageApi.POST("", is.Create)
}

func (is *ImageMetadataServer) Create(gc *gin.Context) {
	ctx, logger := pkgLogger.Ctx(context.Background())
	requestBody := &dtos.CreateImageMetadataRequest{}

	if err := gc.BindJSON(&requestBody); err != nil {
		logger.WithError(err).Error("INVALID_REQUEST")
		middlewares.ErrorResponse(gc, ctx, err)
		return
	}

	if err := requestBody.Validate(); err != nil {
		logger.WithError(err).Error("VALIDATION_FAILURE")
		middlewares.ErrorResponse(gc, ctx, err)
		return
	}

	image, err := is.service.CreateImageMetadata(ctx, requestBody)
	if err != nil {
		middlewares.ErrorResponse(gc, ctx, err)
		return
	}

	response := dtos.ImageMetadataResponseFromModel(image)
	gc.JSON(http.StatusOK, response)
}
