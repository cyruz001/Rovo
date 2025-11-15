package repository

import (
	"context"
	"errors"

	"goServer/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create creates a new notification
func (r *NotificationRepository) Create(ctx context.Context, n *model.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

// FindByID finds a notification by ID
func (r *NotificationRepository) FindByID(ctx context.Context, id string) (*model.Notification, error) {
	var n model.Notification
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&n).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

// Update updates a notification
func (r *NotificationRepository) Update(ctx context.Context, n *model.Notification) error {
	return r.db.WithContext(ctx).Model(n).Updates(n).Error
}

// Delete deletes a notification
func (r *NotificationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Notification{}).Error
}

// GetByUserID gets all notifications for a user
func (r *NotificationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]model.Notification, error) {
	var notifications []model.Notification
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

// GetUnreadCount gets count of unread notifications for a user
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND read = false", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ?", id).
		Update("read", true).Error
}

// MarkAllAsRead marks all notifications as read for a user
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Update("read", true).Error
}

// DeleteByUserID deletes all notifications for a user
func (r *NotificationRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.Notification{}).Error
}

// GetByType gets notifications of a specific type for a user
func (r *NotificationRepository) GetByType(ctx context.Context, userID, notificationType string, limit, offset int) ([]model.Notification, error) {
	var notifications []model.Notification
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, notificationType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
