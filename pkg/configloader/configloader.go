package configloader

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

// Default options for configuration loading.
const (
	DefaultConfigType     = "toml"
	DefaultConfigDir      = "./config"
	DefaultConfigFileName = "default"
	KeyWorkDirEnv         = "WORKDIR"
	// AppModeKey is the key set as
	// environment variable on the kube-manifests deployment yaml
	// env:
	//    - name: APP_MODE
	//      value: test
	AppModeKey = "APP_MODE"
	// AppModeTest means the app is running in test mode
	// In test mode, the database is different
	// Based on the mode, we pick a different configuration file
	AppModeTest = "test"
)

// Options is config options.
type Options struct {
	configType            string
	configPath            string
	testMode              bool
	defaultConfigFileName string
}

// Loader a wrapper over a underlying config loader implementation.
type Loader struct {
	opts  Options
	viper *viper.Viper
}

// NewDefaultOptions returns default options.
// DISCLAIMER: This function is a bit hacky
// This function expects an env $WORKDIR to
// be set and reads configs from $WORKDIR/configs.
// If $WORKDIR is not set. It uses the absolute path wrt
// the location of this file (config.go) to set configPath
// to 2 levels up in viper (../../configs).
// This function breaks if :
// 1. $WORKDIR is set and configs dir not present in $WORKDIR
// 2. $WORKDIR is not set and ../../configs is not present
// 3. $WORKDIR is not set and runtime absolute path of configs
// is different than build time path as runtime.Caller() evaluates
// only at build time
func NewDefaultOptions() Options {
	var configDir string

	workDir := os.Getenv(KeyWorkDirEnv)
	if workDir != "" {
		// used in containers:
		// expects the variable to be set in env
		// $WORKDIR/config
		configDir = path.Join(workDir, DefaultConfigDir)
	} else {
		// used in:
		// $ go run cmd/server/main.go
		// $ bin/darwin_amd64/server
		// ./config
		configDir = DefaultConfigDir
	}

	return NewOptions(DefaultConfigType, configDir, DefaultConfigFileName)
}

// NewOptions returns new Options struct.
func NewOptions(
	configType string,
	configPath string,
	defaultConfigFileName string) Options {

	testMode := false
	if os.Getenv(AppModeKey) == AppModeTest {
		testMode = true
	}

	return Options{configType, configPath, testMode, defaultConfigFileName}
}

// NewDefaultLoader returns new config struct with default options.
func NewDefaultLoader() *Loader {
	return NewLoader(NewDefaultOptions())
}

// NewLoader returns new config struct.
func NewLoader(opts Options) *Loader {
	return &Loader{opts, viper.New()}
}

// Load reads environment specific configurations and along with the defaults
// unmarshalls into config.
func (c *Loader) Load(env string, config interface{}) error {
	// loads the default file then override it with the env file
	err := c.loadByConfigName(c.opts.defaultConfigFileName, config)
	if err != nil {
		return err
	}

	if c.opts.testMode {
		env = env + "_test"
	}

	return c.loadByConfigName(env, config)
}

// loadByConfigName reads configuration from file and unmarshalls into config.
func (c *Loader) loadByConfigName(configName string, config interface{}) error {
	if configName == DefaultConfigFileName {
		fmt.Printf(
			"Loading default config file: %v/%v.%v\n",
			c.opts.configPath,
			configName,
			c.opts.configType,
		)
	} else {
		fmt.Printf(
			"Loading config file: %v/%v.%v\n",
			c.opts.configPath,
			configName,
			c.opts.configType,
		)
	}

	c.viper.SetConfigName(configName)
	c.viper.SetConfigType(c.opts.configType)
	c.viper.AddConfigPath(c.opts.configPath)
	c.viper.AutomaticEnv()
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	return c.viper.Unmarshal(config)
}
