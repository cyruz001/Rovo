package middleware

import "github.com/gofiber/fiber/v3"

func RequireRole(roles ...string) fiber.Handler {
	roleMap := make(map[string]bool)
	for _, role := range roles {
		roleMap[role] = true
	}

	return func(c fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		data := user.(fiber.Map)
		userRole := data["role"].(string)

		if !roleMap[userRole] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden",
			})
		}

		return c.Next()
	}
}
