package api

import (
	"goServer/internal/config"
	"goServer/internal/db"
	"goServer/internal/router"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.Load()

	pool := db.Connect(cfg)
	defer pool.Close()

	app := fiber.New()

	router.SetupRoutes(app, pool)

	log.Println("Server running on", cfg.APP_PORT)
	app.Listen(cfg.APP_PORT)
}
