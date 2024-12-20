package image_metadata

import (
	"fmt"
	"net/http"
	"time"

	internaErr "github.com/danushk97/image-analyzer/internal/errors"
	"github.com/danushk97/image-analyzer/internal/image_metadata/dtos"
	"github.com/danushk97/image-analyzer/internal/image_metadata/model/v1"
	"github.com/danushk97/image-analyzer/internal/image_metadata/service"
	"github.com/danushk97/image-analyzer/internal/middlewares"
	"github.com/danushk97/image-analyzer/pkg/errors"
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
	imageApi.Use(middlewares.AuthMiddleware())

	imageApi.POST("", is.Create)
}

func (is *ImageMetadataServer) Create(gc *gin.Context) {
	var err errors.IError // This will be captured by the defer function
	fn := is.trackRequest(gc)
	defer func() {
		fn(err) // The deferred function uses 'err'
	}()
	logger := pkgLogger.Ctx(gc.Request.Context())

	requestBody := &dtos.CreateImageMetadataRequest{}

	// Bind JSON and assign to 'err'
	if ierr := gc.BindJSON(&requestBody); err != nil {
		err = errors.NewServerError(internaErr.ServerError)
		logger.WithError(err).Error("INVALID_REQUEST")
		middlewares.ErrorResponse(gc, ierr)
		return
	}

	// Validate request body
	if err = requestBody.Validate(); err != nil {
		logger.WithError(err).Error("VALIDATION_FAILURE")
		middlewares.ErrorResponse(gc, err)
		return
	}

	// Create image metadata
	var image *model.ImageMetadata
	image, err = is.service.CreateImageMetadata(
		gc.Request.Context(),
		requestBody,
	)
	if err != nil {
		middlewares.ErrorResponse(gc, err)
		return
	}

	// Send successful response
	response := dtos.ImageMetadataResponseFromModel(image)
	gc.JSON(http.StatusOK, response)
}

// / that logs the latency and final status (success or failure).
func (is *ImageMetadataServer) trackRequest(
	gc *gin.Context,
) func(err errors.IError) {
	logger := pkgLogger.Ctx(gc.Request.Context())

	// Capture the start time
	startTime := time.Now()

	// Log when action starts
	logger.Info("ACTION_STARTED")

	// Return a defer function to log final status and latency
	return func(err errors.IError) {
		latency := time.Since(startTime) // Calculate latency
		if err != nil {
			logger.WithError(err).Error(
				fmt.Sprintf("ACTION_FAILED | latency: %v", latency),
			)
		} else {
			logger.Info(
				fmt.Sprintf("ACTION_SUCCESS | latency: %v", latency),
			)
		}
	}
}
