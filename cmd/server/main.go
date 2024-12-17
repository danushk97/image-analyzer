package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/danushk97/image-analyzer/internal/config"
	health "github.com/danushk97/image-analyzer/internal/health"
	"github.com/danushk97/image-analyzer/internal/image_metadata"
	imageMetaCore "github.com/danushk97/image-analyzer/internal/image_metadata/service"
	srv "github.com/danushk97/image-analyzer/internal/server"
	"github.com/danushk97/image-analyzer/pkg/env"
	pkgLogger "github.com/danushk97/image-analyzer/pkg/logger"
	"github.com/danushk97/image-analyzer/pkg/storage"
)

func main() {
	env := env.GetEnv()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := pkgLogger.NewLogger()

	// load configurations and distribute parts of it in main
	config := config.NewConfig(env)

	// storage service is the service for main persistent store
	storageService, err := storage.New(ctx, config.Store)
	if err != nil {
		logger.Fatalf(
			"could not create database, err:%+v", err,
		)
	}

	imageMetaService := imageMetaCore.NewService(
		imageMetaCore.WithStorage(storageService),
	)

	healthServer := health.NewServer()

	imageServer := image_metadata.NewServer(imageMetaService)

	server := srv.New(ctx, &srv.Config{})

	server.WithOptions(
		server.WithHealthServer(healthServer),
		server.WithImageMetadataServer(imageServer),
	)

	// graceful shutdown, no libs required, understand just below
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigterm:
			logger.Info(ctx, "sigterm received")
		case <-ctx.Done():
			logger.Info(ctx, "context done, bye")
			return
		}

		err := server.Shutdown(ctx)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("error shutting down server(s):%+v", err))
		}

		logger.Info(ctx, "cancel() context")
		cancel()
	}()

	// run starts http, grpc, and metric server
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = server.Run(ctx)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("server shutdown with err(s):%+v", err))
			cancel()
		}
	}()
	logger.Info(ctx, "server(s) running", "log_level", logger.Level())

	// wait for all go routines to shutdown, then exit main
	wg.Wait()

	logger.Info(ctx, "gracefully shutdown")
}
