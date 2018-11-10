package infra

import (
	"log"
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
	log.Printf("Environment is '%s'.\n", getEnvironment())

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

func isProduction() bool {
	return getEnvironment() == prod
}
