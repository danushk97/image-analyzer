package health

import (
	"context"
	"net/http"

	pkgLogger "github.com/danushk97/image-analyzer/pkg/logger"
	"github.com/gin-gonic/gin"
)

type HealthServer struct {
}

func NewServer() *HealthServer {
	return &HealthServer{}
}

func (hs *HealthServer) SetupRoutes(r *gin.Engine) {
	healthApi := r.Group("/v1/health")

	healthApi.GET("", hs.Check)
}

func (hs *HealthServer) Check(gc *gin.Context) {
	logger := pkgLogger.Ctx(context.Background())

	logger.Info("Health check.")

	gc.JSON(http.StatusOK, gin.H{"status": "up"})
}
