package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	"goServer/internal/config"
	"goServer/internal/db"
	"goServer/internal/router"
)

func init() {
	// Try local folder (cmd/api)
	if err := godotenv.Load(); err != nil {
		// Try project root
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("[init] No .env file found in current or parent dir")
		}
	}
}

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	database := db.Connect(cfg)

	app := fiber.New()

	router.SetupRoutes(app, database, cfg)

	log.Printf("Server listening on port %s", cfg.AppPort)
	if err := app.Listen(cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
