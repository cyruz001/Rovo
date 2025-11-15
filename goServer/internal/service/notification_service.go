package service

import (
	"context"
	"errors"
	"fmt"

	"goServer/internal/model"
	"goServer/internal/repository"
)

type NotificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

func NewNotificationService(nr repository.NotificationRepository, ur repository.UserRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: nr,
		userRepo:         ur,
	}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(ctx context.Context, userID, notificationType string) (*model.Notification, error) {
	if userID == "" || notificationType == "" {
		return nil, errors.New("user id and notification type are required")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	notification := &model.Notification{
		UserID: userID,
		Type:   notificationType,
		Read:   false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

// GetNotifications retrieves user's notifications with pagination
func (s *NotificationService) GetNotifications(ctx context.Context, userID string, limit, offset int) ([]model.Notification, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	notifications, err := s.notificationRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}

// GetUnreadCount gets count of unread notifications
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, errors.New("user id is required")
	}

	count, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID string) (*model.Notification, error) {
	if notificationID == "" || userID == "" {
		return nil, errors.New("notification id and user id are required")
	}

	notification, err := s.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}
	if notification == nil {
		return nil, errors.New("notification not found")
	}

	if notification.UserID != userID {
		return nil, errors.New("unauthorized: can only read your own notifications")
	}

	notification.Read = true

	if err := s.notificationRepo.Update(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to update notification: %w", err)
	}

	return notification, nil
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID, userID string) error {
	if notificationID == "" || userID == "" {
		return errors.New("notification id and user id are required")
	}

	notification, err := s.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("failed to find notification: %w", err)
	}
	if notification == nil {
		return errors.New("notification not found")
	}

	if notification.UserID != userID {
		return errors.New("unauthorized: can only delete your own notifications")
	}

	if err := s.notificationRepo.Delete(ctx, notificationID); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user id is required")
	}

	if err := s.notificationRepo.MarkAllAsRead(ctx, userID); err != nil {
		return fmt.Errorf("failed to mark all as read: %w", err)
	}

	return nil
}

// GetNotificationsByType retrieves notifications of a specific type
func (s *NotificationService) GetNotificationsByType(ctx context.Context, userID, notificationType string, limit, offset int) ([]model.Notification, error) {
	if userID == "" || notificationType == "" {
		return nil, errors.New("user id and notification type are required")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	notifications, err := s.notificationRepo.GetByType(ctx, userID, notificationType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications by type: %w", err)
	}

	return notifications, nil
}

// NotifyPostLike creates a notification when someone likes a post
func (s *NotificationService) NotifyPostLike(ctx context.Context, postOwnerID, likerID string) error {
	if postOwnerID == "" || likerID == "" {
		return errors.New("post owner id and liker id are required")
	}

	if postOwnerID == likerID {
		return nil // Don't notify on self-like
	}

	notification := &model.Notification{
		UserID: postOwnerID,
		Type:   "LIKE",
		Read:   false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create like notification: %w", err)
	}

	return nil
}

// NotifyPostRepost creates a notification when someone reposts a post
func (s *NotificationService) NotifyPostRepost(ctx context.Context, postOwnerID, reposterID string) error {
	if postOwnerID == "" || reposterID == "" {
		return errors.New("post owner id and reposter id are required")
	}

	if postOwnerID == reposterID {
		return nil // Don't notify on self-repost
	}

	notification := &model.Notification{
		UserID: postOwnerID,
		Type:   "REPOST",
		Read:   false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create repost notification: %w", err)
	}

	return nil
}

// NotifyPostReply creates a notification when someone replies to a post
func (s *NotificationService) NotifyPostReply(ctx context.Context, postOwnerID, replierID string) error {
	if postOwnerID == "" || replierID == "" {
		return errors.New("post owner id and replier id are required")
	}

	if postOwnerID == replierID {
		return nil // Don't notify on self-reply
	}

	notification := &model.Notification{
		UserID: postOwnerID,
		Type:   "REPLY",
		Read:   false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create reply notification: %w", err)
	}

	return nil
}

// NotifyMention creates a notification when someone mentions a user
func (s *NotificationService) NotifyMention(ctx context.Context, mentionedUserID, mentionerID string) error {
	if mentionedUserID == "" || mentionerID == "" {
		return errors.New("mentioned user id and mentioner id are required")
	}

	if mentionedUserID == mentionerID {
		return nil // Don't notify on self-mention
	}

	notification := &model.Notification{
		UserID: mentionedUserID,
		Type:   "MENTION",
		Read:   false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create mention notification: %w", err)
	}

	return nil
}

// NotifyFollow creates a notification when someone follows a user
func (s *NotificationService) NotifyFollow(ctx context.Context, followeeID, followerID string) error {
	if followeeID == "" || followerID == "" {
		return errors.New("followee id and follower id are required")
	}

	notification := &model.Notification{
		UserID: followeeID,
		Type:   "FOLLOW",
		Read:   false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create follow notification: %w", err)
	}

	return nil
}

// DeleteNotificationsByUserID deletes all notifications for a user
func (s *NotificationService) DeleteNotificationsByUserID(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user id is required")
	}

	if err := s.notificationRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete notifications: %w", err)
	}

	return nil
}
