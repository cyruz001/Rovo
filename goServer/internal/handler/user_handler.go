package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"goServer/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Profile(c fiber.Ctx) error {
	// retrieve token object stored by middleware
	raw := c.Locals("user")
	if raw == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	token, ok := raw.(*jwt.Token)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token structure")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
	}

	id, ok := claims["sub"].(string)
	if !ok || id == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing sub claim")
	}

	user, err := h.service.GetByID(context.Background(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "database error")
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(user)
}
