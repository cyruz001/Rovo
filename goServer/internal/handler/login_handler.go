package handler

import (
	"time"

	"goServer/internal/dto"
	"goServer/internal/model"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginHandler(db *gorm.DB, jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var body dto.LoginRequest
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}

		var user model.User
		if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid email"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"role":    user.Role,
			"exp":     time.Now().Add(time.Hour * 1).Unix(), // 1 hour
		})

		tokenString, _ := token.SignedString([]byte(jwtSecret))

		return c.JSON(fiber.Map{
			"message": "Login successful",
			"token":   tokenString,
		})
	}
}
