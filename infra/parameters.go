package infra

import (
	"fmt"
	"log"
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
	SlackClientID        string `env:"SLACK_CLIENT_ID,required"`
	SlackClientSecret    string `env:"SLACK_CLIENT_SECRET,required"`
	SlackSignSecret      string `env:"SLACK_SIGN_SECRET,required"`
	SlackRedirectedPath  string `env:"SLACK_REDIRECTED_PATH,required"`
	APIHost              string `env:"API_HOST,required"`
	GoogleClientID       string `env:"GOOGLE_CLIENT_ID,required"`
	GoogleClientSecret   string `env:"GOOGLE_CLIENT_SECRET,required"`
	GoogleRedirectedPath string `env:"GOOGLE_REDIRECTED_PATH,required"`
	LogWebhook           string `env:"LOG_WEBHOOK"`
}

const (
	repoPath  = "./"
	envPrefix = ".slgmails"
)

func setupEnv() {
	log.Println("Setup environment variable...")
	loadEnv()
	parseEnv()
}

func loadEnv() {
	rootPath := filepath.Join(repoPath, envFileName())
	err := godotenv.Load(rootPath)
	if err != nil {
		log.Fatal("loading env file error", err)
	}
}

func parseEnv() {
	Env = &envConfig{}
	err := env.Parse(Env)
	if err != nil {
		log.Fatal("parse env file error", err)
	}
}

func envFileName() string {
	return fmt.Sprintf("%s.%s.env", envPrefix, getEnvironment())
}
