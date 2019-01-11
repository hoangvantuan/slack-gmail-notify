package infra

import (
	"fmt"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Env is enviroment variable
var Env *envConfig

type envConfig struct {
	EncryptKey           string `env:"SLGMAILS_ENCRYPT_KEY,required"`
	MysqlUser            string `env:"SLGMAILS_MYSQL_MASTER_USENAME,required"`
	MysqlPass            string `env:"SLGMAILS_MYSQL_MASTER_PASSWORD,required"`
	MysqlEndpoint        string `env:"SLGMAILS_MYSQL_ENDPOINT,required"`
	Port                 string `env:"SLGMAILS_PORT" envDefault:"8080"`
	SlackClientID        string `env:"SLACK_CLIENT_ID"`
	SlackClientSecret    string `env:"SLACK_CLIENT_SECRET"`
	SlackSignSecret      string `env:"SLACK_SIGN_SECRET"`
	SlackRedirectedPath  string `env:"SLACK_REDIRECTED_PATH"`
	APIHost              string `env:"API_HOST"`
	GoogleClientID       string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret   string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectedPath string `env:"GOOGLE_REDIRECTED_PATH"`
}

const (
	repoPath = "./"
	envExt   = ".env"
)

func setupEnv() {
	loadEnv()
	parseEnv()
}

func loadEnv() {
	rootPath := filepath.Join(repoPath, envFileName())
	err := godotenv.Load(rootPath)
	if err != nil {
		Fatal("loading env file error", err)
	}
}

func parseEnv() {
	Env = &envConfig{}
	err := env.Parse(Env)
	if err != nil {
		Fatal("parse env file error", err)
	}
}

func envFileName() string {
	return fmt.Sprintf("%s.%s", envExt, getEnvironment())
}
