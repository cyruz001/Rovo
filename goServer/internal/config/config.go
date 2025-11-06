package config

import (
	"os"
)

type Config struct {
	APP_PORT string
	DB_URL   string
}

func Load() Config {
	return Config{
		APP_PORT: getEnv("APP_PORT", ":3000"),
		DB_URL:   getEnv("DB_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
