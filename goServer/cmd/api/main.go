package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"

	"goServer/internal/config"
	"goServer/internal/db"
	"goServer/internal/router"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	pool := db.Connect(cfg)
	defer func() {
		sqlDB, _ := pool.DB()
		_ = sqlDB.Close()
	}()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	router.SetupRoutes(app, pool, cfg)

	log.Printf("Server listening on %s", cfg.AppPort)
	if err := app.Listen(cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
