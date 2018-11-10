package infra

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Env is enviroment variable
var Env *envConfig

type envConfig struct {
	EncryptKey    string `env:"SLGMAILS_ENCRYPT_KEY,required"`
	MysqlUser     string `env:"SLGMAILS_MYSQL_MASTER_USENAME,required"`
	MysqlPass     string `env:"SLGMAILS_MYSQL_MASTER_PASSWORD,required"`
	MysqlEndpoint string `env:"SLGMAILS_MYSQL_ENDPOINT" envDefault:"ec2-52-193-31-72.ap-northeast-1.compute.amazonaws.com:3306"`
	Port          string `env:"SLGMAILS_PORT" envDefault:"8080"`
}

func setupEnv() {
	loadEnv()
	parseEnv()
}

func loadEnv() {
	envFileName := fmt.Sprintf(".env.%s", getEnvironment())
	rootPath := filepath.Join(os.Getenv("GOPATH"), "src/github.com/mdshun/slack-gmail-notify", envFileName)
	err := godotenv.Load(rootPath)

	if err != nil {
		Swarn("loading .env.%s file\n%s", getEnvironment(), err)
	}
}

func parseEnv() {
	Env = &envConfig{}

	err := env.Parse(Env)

	if err != nil {
		Swarn("[Warn] parse envs to struct\n%s", err)
	}
}
