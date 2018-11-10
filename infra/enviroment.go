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
	case dev:
		return "dev"
	case stg:
		return "stg"
	case prod:
		return "prod"
	default:
		return "unknown"
	}
}

// Setup prepares database connection and parameters for running application.
func Setup() {
	if RDB != nil {
		Linfo("Infra already setup!")
		return
	}

	setupLogger()

	Sinfo("Environment is ", getEnvironment())

	setupEnv()
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

// IsProduction is true if current env is prod
func IsProduction() bool {
	return getEnvironment() == prod
}
