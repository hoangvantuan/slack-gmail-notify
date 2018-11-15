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
	EncryptKey          string   `env:"SLGMAILS_ENCRYPT_KEY,required"`
	MysqlUser           string   `env:"SLGMAILS_MYSQL_MASTER_USENAME,required"`
	MysqlPass           string   `env:"SLGMAILS_MYSQL_MASTER_PASSWORD,required"`
	MysqlEndpoint       string   `env:"SLGMAILS_MYSQL_ENDPOINT,required"`
	Port                string   `env:"SLGMAILS_PORT" envDefault:"8080"`
	SlackClientID       string   `env:"SLACK_CLIENT_ID"`
	SlackClientSecret   string   `env:"SLACK_CLIENT_SECRET"`
	SlackSignSecret     string   `env:"SLACK_SIGN_SECRET"`
	SlackRedirectedURL  string   `env:"SLACK_REDIRECTED_URL"`
	SlackScope          string   `env:"SLACK_SCOPES"`
	APIHost             string   `env:"API_HOST"`
	GoogleClientID      string   `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret  string   `env:"GOOGLE_CLIENT_SECRET"`
	GoogleAuthURL       string   `env:"GOOGLE_AUTH_URL"`
	GoogleTokenURL      string   `env:"GOOGLE_TOKEN_URL"`
	GoogleScopes        []string `env:"GOOGLE_SCOPES"`
	GoogleRedirectedURL string   `env:"GOOGLE_REDIRECTED_URL"`
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
		Swarn("parse envs to struct\n%s", err)
	}
}
