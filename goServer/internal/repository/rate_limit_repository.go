package repository

import (
	"context"
	"errors"
	"time"

	"goServer/internal/model"

	"gorm.io/gorm"
)

type RateLimitRepository struct {
	db *gorm.DB
}

func NewRateLimitRepository(db *gorm.DB) *RateLimitRepository {
	return &RateLimitRepository{db: db}
}

// Create creates a new rate limit record
func (r *RateLimitRepository) Create(ctx context.Context, rl *model.RateLimit) error {
	return r.db.WithContext(ctx).Create(rl).Error
}

// FindByUserAndAction finds a rate limit record
func (r *RateLimitRepository) FindByUserAndAction(ctx context.Context, userID, action string) (*model.RateLimit, error) {
	var rl model.RateLimit
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND action = ?", userID, action).
		First(&rl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rl, nil
}

// GetCount gets the count for a user action in a time window
func (r *RateLimitRepository) GetCount(ctx context.Context, userID, action string, windowStart time.Time) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.RateLimit{}).
		Where("user_id = ? AND action = ? AND window_start >= ?", userID, action, windowStart).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// IncrementCount increments the count for a user action
func (r *RateLimitRepository) IncrementCount(ctx context.Context, userID, action string, windowStart time.Time) error {
	rl := &model.RateLimit{
		UserID:      userID,
		Action:      action,
		WindowStart: windowStart,
		Count:       1,
	}

	return r.db.WithContext(ctx).
		Model(rl).
		Where("user_id = ? AND action = ? AND window_start = ?", userID, action, windowStart).
		Updates(map[string]interface{}{
			"count": gorm.Expr("count + 1"),
		}).Error
}

// Reset resets the rate limit for a user action
func (r *RateLimitRepository) Reset(ctx context.Context, userID, action string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND action = ?", userID, action).
		Delete(&model.RateLimit{}).Error
}

// DeleteExpired deletes expired rate limit records
func (r *RateLimitRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("window_start < ?", before).
		Delete(&model.RateLimit{}).Error
}
