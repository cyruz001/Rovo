package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"goServer/internal/model"
	"goServer/internal/repository"
)

type RateLimitService struct {
	rateLimitRepo repository.RateLimitRepository
}

func NewRateLimitService(rlr repository.RateLimitRepository) *RateLimitService {
	return &RateLimitService{rateLimitRepo: rlr}
}

// CheckLimit checks if a user has exceeded rate limit
func (s *RateLimitService) CheckLimit(ctx context.Context, userID, action string, limit int, window time.Duration) (bool, error) {
	if userID == "" || action == "" {
		return false, errors.New("user id and action are required")
	}

	now := time.Now()
	windowStart := now.Add(-window)

	// Get current count
	count, err := s.rateLimitRepo.GetCount(ctx, userID, action, windowStart)
	if err != nil {
		return false, fmt.Errorf("failed to get rate limit count: %w", err)
	}

	// If count exceeds limit, deny
	if count >= limit {
		return false, nil
	}

	// Increment count
	if err := s.rateLimitRepo.IncrementCount(ctx, userID, action, now); err != nil {
		return false, fmt.Errorf("failed to increment rate limit: %w", err)
	}

	return true, nil
}

// GetRateLimit retrieves rate limit info
func (s *RateLimitService) GetRateLimit(ctx context.Context, userID, action string) (*model.RateLimit, error) {
	if userID == "" || action == "" {
		return nil, errors.New("user id and action are required")
	}

	rl, err := s.rateLimitRepo.FindByUserAndAction(ctx, userID, action)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit: %w", err)
	}

	return rl, nil
}

// ResetRateLimit resets rate limit for a user action
func (s *RateLimitService) ResetRateLimit(ctx context.Context, userID, action string) error {
	if userID == "" || action == "" {
		return errors.New("user id and action are required")
	}

	if err := s.rateLimitRepo.Reset(ctx, userID, action); err != nil {
		return fmt.Errorf("failed to reset rate limit: %w", err)
	}

	return nil
}
