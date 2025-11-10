package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	AppPort     string
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() Config {
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		AppPort:     os.Getenv("APP_PORT"),
	}
}
