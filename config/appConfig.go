package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort   string
	DSN          string
	MigrationURL string
	DBSource     string
}

func SetupEnv() (cfg AppConfig, err error) {
	if os.Getenv("APP_ENV") == "dev" {
		godotenv.Load()
	}

	httpPort := os.Getenv("HTTP_PORT")
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("HTTP_PORT is not set")
	}

	dsn := os.Getenv("DSN")
	MigrationURL := os.Getenv("MIGRATION_URL")
	DBSource := os.Getenv("DB_SOURCE")

	if len(dsn) < 1 {
		return AppConfig{}, errors.New("DSN is not set")
	}
	return AppConfig{
		ServerPort:   httpPort,
		DSN:          dsn,
		MigrationURL: MigrationURL,
		DBSource:     DBSource,
	}, nil
}
