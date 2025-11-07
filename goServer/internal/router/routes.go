package router

import (
	"github.com/gofiber/fiber/v3"

	"goServer/internal/config"
	"goServer/internal/handler"
	"goServer/internal/repository"
	"goServer/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, cfg config.Config) {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(userSvc, cfg)
	userHandler := handler.NewUserHandler(userSvc)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/auth/register", authHandler.Register)
	v1.Post("/auth/login", authHandler.Login)

	// Protected group
	protected := v1.Group("/user")
	protected.Use(func(c fiber.Ctx) error {
		auth := c.Get("Authorization")
		if len(auth) < 7 || auth[:7] != "Bearer " {
			return fiber.NewError(fiber.StatusUnauthorized, "missing or invalid token")
		}

		tokenStr := auth[7:]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		c.Locals("user", token)
		return c.Next()
	})

	protected.Get("/me", userHandler.Profile)
}
