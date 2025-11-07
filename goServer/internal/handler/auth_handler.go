package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"goServer/internal/config"
	"goServer/internal/service"
)

type AuthHandler struct {
	service *service.UserService
	cfg     config.Config
}

func NewAuthHandler(s *service.UserService, cfg config.Config) *AuthHandler {
	return &AuthHandler{service: s, cfg: cfg}
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req registerReq

	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	user, err := h.service.Register(context.Background(), req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	})
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq

	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	user, err := h.service.Authenticate(context.Background(), req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	// JWT
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to sign token")
	}

	return c.JSON(fiber.Map{"access_token": signed})
}
