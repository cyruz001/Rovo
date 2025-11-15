package handler

import (
	"context"
	"strconv"

	"goServer/internal/dto"
	"goServer/internal/model"
	"goServer/internal/service"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService         *service.UserService
	notificationService *service.NotificationService
}

func NewUserHandler(us *service.UserService, ns *service.NotificationService) *UserHandler {
	return &UserHandler{
		userService:         us,
		notificationService: ns,
	}
}

// GetProfile retrieves current user profile
func (h *UserHandler) GetProfile(c fiber.Ctx) error {
	userID, ok := c.Locals("sub").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, err := h.userService.GetUserByID(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(userToRes(user))
}

// UpdateProfile updates current user profile
func (h *UserHandler) UpdateProfile(c fiber.Ctx) error {
	userID, ok := c.Locals("sub").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req dto.UserUpdateReq

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := h.userService.UpdateUser(context.Background(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(userToRes(user))
}

// DeleteAccount deletes current user account
func (h *UserHandler) DeleteAccount(c fiber.Ctx) error {
	userID, ok := c.Locals("sub").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.userService.DeleteUser(context.Background(), userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "account deleted successfully"})
}

// GetUserByUsername retrieves user by username
func (h *UserHandler) GetUserByUsername(c fiber.Ctx) error {
	username := c.Params("username")

	user, err := h.userService.GetUserByUsername(context.Background(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(userToRes(user))
}

// GetFollowers retrieves followers of a user
func (h *UserHandler) GetFollowers(c fiber.Ctx) error {
	username := c.Params("username")

	followers, err := h.userService.GetFollowers(context.Background(), username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.UserRes, len(followers))
	for i, f := range followers {
		res[i] = userToRes(&f)
	}

	return c.JSON(res)
}

// GetFollowing retrieves users that a user is following
func (h *UserHandler) GetFollowing(c fiber.Ctx) error {
	username := c.Params("username")

	following, err := h.userService.GetFollowing(context.Background(), username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.UserRes, len(following))
	for i, f := range following {
		res[i] = userToRes(&f)
	}

	return c.JSON(res)
}

// GetMyFollowers retrieves current user's followers
func (h *UserHandler) GetMyFollowers(c fiber.Ctx) error {
	userID, ok := c.Locals("sub").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, err := h.userService.GetUserByID(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	followers, err := h.userService.GetFollowers(context.Background(), user.Username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.UserRes, len(followers))
	for i, f := range followers {
		res[i] = userToRes(&f)
	}

	return c.JSON(res)
}

// GetMyFollowing retrieves users that current user is following
func (h *UserHandler) GetMyFollowing(c fiber.Ctx) error {
	userID, ok := c.Locals("sub").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, err := h.userService.GetUserByID(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	following, err := h.userService.GetFollowing(context.Background(), user.Username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.UserRes, len(following))
	for i, f := range following {
		res[i] = userToRes(&f)
	}

	return c.JSON(res)
}

// FollowUser follows a user
func (h *UserHandler) FollowUser(c fiber.Ctx) error {
	followerID, ok := c.Locals("sub").(string)
	if !ok || followerID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	followeeID := c.Params("id")

	if err := h.userService.FollowUser(context.Background(), followerID, followeeID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "followed successfully"})
}

// UnfollowUser unfollows a user
func (h *UserHandler) UnfollowUser(c fiber.Ctx) error {
	followerID, ok := c.Locals("sub").(string)
	if !ok || followerID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	followeeID := c.Params("id")

	if err := h.userService.UnfollowUser(context.Background(), followerID, followeeID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "unfollowed successfully"})
}

// SearchUsers searches for users
func (h *UserHandler) SearchUsers(c fiber.Ctx) error {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query is required"})
	}

	users, err := h.userService.SearchUsers(context.Background(), query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.UserRes, len(users))
	for i, u := range users {
		res[i] = userToRes(&u)
	}

	return c.JSON(res)
}

// GetAllUsers retrieves all users (admin)
func (h *UserHandler) GetAllUsers(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, err := h.userService.GetAllUsers(context.Background(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.UserRes, len(users))
	for i, u := range users {
		res[i] = userToRes(&u)
	}

	return c.JSON(res)
}

// AdminDeleteUser deletes a user (admin)
func (h *UserHandler) AdminDeleteUser(c fiber.Ctx) error {
	userID := c.Params("id")

	if err := h.userService.DeleteUser(context.Background(), userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "user deleted successfully"})
}

// UpdateUserRole updates a user's role (admin)
func (h *UserHandler) UpdateUserRole(c fiber.Ctx) error {
	userID := c.Params("id")
	var req dto.UpdateUserRoleReq

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := h.userService.UpdateUserRole(context.Background(), userID, req.Role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(userToRes(user))
}

// GetNotifications retrieves user's notifications
func (h *UserHandler) GetNotifications(c fiber.Ctx) error {
	userID, ok := c.Locals("sub").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	notifications, err := h.notificationService.GetNotifications(context.Background(), userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.NotificationRes, len(notifications))
	for i, n := range notifications {
		res[i] = notificationToRes(&n)
	}

	return c.JSON(res)
}

// MarkNotificationAsRead marks a notification as read
func (h *UserHandler) MarkNotificationAsRead(c fiber.Ctx) error {
	userID, ok := c.Locals("sub").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	notificationID := c.Params("id")

	notification, err := h.notificationService.MarkAsRead(context.Background(), notificationID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(notificationToRes(notification))
}

// DeleteNotification deletes a notification
func (h *UserHandler) DeleteNotification(c fiber.Ctx) error {
	userID, ok := c.Locals("sub").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	notificationID := c.Params("id")

	if err := h.notificationService.DeleteNotification(context.Background(), notificationID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "notification deleted"})
}

// GetSystemStats retrieves system statistics (admin)
func (h *UserHandler) GetSystemStats(c fiber.Ctx) error {
	// This will be implemented when stats service is created
	return c.JSON(fiber.Map{"message": "system stats"})
}

// Helper function to convert User model to UserRes DTO
func userToRes(u *model.User) dto.UserRes {
	return dto.UserRes{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		AvatarURL:   u.AvatarURL,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt.String(),
	}
}

// Helper function to convert Notification model to NotificationRes DTO
func notificationToRes(n *model.Notification) dto.NotificationRes {
	return dto.NotificationRes{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Message:   n.Type,
		Read:      n.Read,
		CreatedAt: n.CreatedAt.String(),
	}
}
