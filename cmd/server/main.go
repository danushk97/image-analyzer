package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/danushk97/image-analyzer/internal/config"
	health "github.com/danushk97/image-analyzer/internal/health"
	srv "github.com/danushk97/image-analyzer/internal/server"
	pkgLogger "github.com/danushk97/image-analyzer/pkg/logger"
	"github.com/danushk97/image-analyzer/pkg/storage"
)

// EnvDev signifies it is a dev and is
// used to decide seed or not seed data
const EnvDev = "dev"

// ModeTest signifies it is a testing mode
const ModeTest = "test"

// ModeLive signifies it is a testing mode
const ModeLive = "live"

// New fetches env for bootstrapping
func getEnv() string {
	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = EnvDev
	}

	return environment
}

func main() {
	env := getEnv()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := pkgLogger.NewLogger()

	// load configurations and distribute parts of it in main
	config := config.NewConfig(env)

	// storage service is the service for main persistent store
	_, err := storage.New(ctx, config.Store)
	if err != nil {
		log.Fatal("could not create database, err:%+v", err)
	}
	healthServer := health.NewHealthServer()

	server := srv.New(ctx, &srv.Config{})

	server.WithOptions(
		server.WithHealthServer(healthServer),
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
