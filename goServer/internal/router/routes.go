package router

import (
	"goServer/internal/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupRoutes(app *fiber.App, dbPool *pgxpool.Pool) {
	userHandler := handler.NewUserHandler(dbPool)
	app.Get("/users/:id", userHandler.GetUser)

}
