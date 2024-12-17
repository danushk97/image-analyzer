package server

import (
	"context"
	"fmt"
	"net/http"

	healthServer "github.com/danushk97/image-analyzer/internal/health"
	"github.com/danushk97/image-analyzer/internal/image_metadata"
	pkgLogger "github.com/danushk97/image-analyzer/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	// DefaultHTTPAddress for HTTP server
	DefaultHTTPAddress = "0.0.0.0:8081"
	// DefaultShutdownTimeout is the default time allowed to shutdown server
	DefaultShutdownTimeout = 20
)

// Config holds the server configurations
type Config struct {
	ShutdownTimeout int
	ServerAddress   string
}

type ServerOption func(s *Server) error

// Server wraps the gin.Engine
type Server struct {
	server *http.Server
	router *gin.Engine
	config *Config
}

// New creates a new server
func New(ctx context.Context, config *Config) *Server {
	ctx, logger := pkgLogger.Ctx(ctx)

	if config.ServerAddress == "" {
		config.ServerAddress = DefaultHTTPAddress
	}
	if config.ShutdownTimeout <= 0 {
		config.ShutdownTimeout = DefaultShutdownTimeout
	}

	router := gin.Default()

	logger.Info(
		fmt.Sprintf(
			"registered server address http_server: %v",
			config.ServerAddress,
		),
	)

	return &Server{
		router: router,
		config: config,
	}
}

// Run starts the Gin server
func (s *Server) Run(ctx context.Context) error {
	s.server = &http.Server{
		Addr:    s.config.ServerAddress,
		Handler: s.router,
	}

	go func() {

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkgLogger.NewLogger().Error(ctx, "server error: %v", err)
		}
	}()

	<-ctx.Done()
	return s.Shutdown(ctx)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, DefaultShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(shutdownCtx)
}

func (s *Server) WithOptions(opts ...ServerOption) error {
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) WithHealthServer(hs *healthServer.HealthServer) ServerOption {
	return func(s *Server) error {
		hs.SetupRoutes(s.router)
		return nil
	}
}

func (s *Server) WithImageMetadataServer(is *image_metadata.ImageMetadataServer) ServerOption {
	return func(s *Server) error {
		is.SetupRoutes(s.router)
		return nil
	}
}
