package middleware

import (
	"context"
	"time"

	"goServer/internal/service"

	"github.com/gofiber/fiber/v3"
)

func RateLimit(rateLimitService *service.RateLimitService, action string, limit int, window time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("sub")
		if userID == nil {
			return c.Next() // Skip rate limit for unauthenticated requests
		}

		allowed, err := rateLimitService.CheckLimit(context.Background(), userID.(string), action, limit, window)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "rate limit check failed"})
		}

		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "rate limit exceeded"})
		}

		return c.Next()
	}
}
