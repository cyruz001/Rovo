package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"

	"goServer/internal/config"
	"goServer/internal/db"
	"goServer/internal/model"
	"goServer/internal/router"
)

func init() {

	if err := godotenv.Load(); err != nil {
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

	//database migrations
	if err := database.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("[main] Database migrations completed successfully")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	router.SetupRoutes(app, database, cfg)

	log.Printf("Server listening on port %s", cfg.AppPort)
	if err := app.Listen(cfg.AppPort); err != nil {
		log.Fatal(err)
	}

}
