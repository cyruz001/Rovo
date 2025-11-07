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
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://postgres:8eGyXSAYimygW4m@db.jzckycyvzxolhrduashn.supabase.co:5432/postgres"),
		JWTSecret:   getEnv("JWT_SECRET", "replace-this-secret"),
		AppPort:     getEnv("APP_PORT", ":3000"),
	}
}
