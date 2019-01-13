package infra

import (
	"os"
)

type environment int

const (
	dev environment = iota
	stg
	prod
)

func (e environment) String() string {
	switch e {
	case stg:
		return "stg"
	case prod:
		return "prod"
	default:
		return "dev"
	}
}

// Setup prepares database connection and parameters for running application.
func Setup() {
	if RDB != nil {
		return
	}

	setupEnv()
	setupLogger()

	Info("Running on ", getEnvironment())

	setupDatabase()
}

func getEnvironment() environment {
	env := dev
	e := os.Getenv("SLGMAILS_ENV")
	if e == "stg" {
		env = stg
	}
	if e == "prod" {
		env = prod
	}

	return env
}

// IsProduction is true if current env is prod or stg
func IsProduction() bool {
	return getEnvironment() == prod || getEnvironment() == stg
}
