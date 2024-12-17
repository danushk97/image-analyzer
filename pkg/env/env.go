package env

import "os"

// EnvDev signifies it is a dev and is
// used to decide seed or not seed data
const EnvDev = "dev"

// ModeTest signifies it is a testing mode
const ModeTest = "test"

// ModeLive signifies it is a testing mode
const ModeLive = "live"

// New fetches env for bootstrapping
func GetEnv() string {
	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = EnvDev
	}

	return environment
}
