package router

import (
	"goServer/internal/config"
	"goServer/internal/handler"
	"goServer/internal/middleware"
	"goServer/internal/repository"
	"goServer/internal/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, cfg config.Config) {
	// dependency injection
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(*userRepo)

	authHandler := handler.NewAuthHandler(userSvc, cfg)
	userHandler := handler.NewUserHandler(userSvc)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public
	v1.Post("/auth/register", authHandler.Register)
	v1.Post("/auth/login", authHandler.Login)

	// Protected group with JWT guard
	protected := v1.Group("/", middleware.JWT(cfg.JWTSecret))

	protected.Get("/user/me", userHandler.Profile)
}
