package config

import (
	"fmt"
	"os"

	"github.com/danushk97/image-analyzer/pkg/configloader"
	"github.com/danushk97/image-analyzer/pkg/storage"
)

// Config holds the entire configuration for the service
type Config struct {
	// App configurations
	App App

	Store storage.Config
}

// App contains application-specific config values
type App struct {
	// Env the application runs
	Env string
	// ServiceName is the of the application
	ServiceName string
	// HostName is URL of the service
	HostName string
	// StoreChoice  is the storage choice
	StoreChoice string // valid values: sql,dynamodb
}

// NewConfig creates new instance of the
// Config from the config file: <env>.toml
func NewConfig(env string) *Config {
	var config Config

	loader := configloader.NewDefaultLoader()
	err := loader.Load(env, &config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return &config
}
