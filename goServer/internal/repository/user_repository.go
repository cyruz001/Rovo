package repository

import (
	"context"
	"errors"

	"goServer/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	return r.db.WithContext(ctx).Model(u).Updates(u).Error
}

// UpdateProfile updates user profile fields
func (r *UserRepository) UpdateProfile(ctx context.Context, id string, displayName, bio, avatarURL string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.WithContext(ctx).Model(u).Where("id = ?", id).Updates(map[string]interface{}{
		"display_name": displayName,
		"bio":          bio,
		"avatar_url":   avatarURL,
	}).First(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}

// FollowUser creates a follow relationship
func (r *UserRepository) FollowUser(ctx context.Context, followerID, followeeID string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", followeeID).
		Association("Followers").
		Append(&model.User{ID: followerID})
}

// UnfollowUser removes a follow relationship
func (r *UserRepository) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", followeeID).
		Association("Followers").
		Delete(&model.User{ID: followerID})
}

// GetFollowers gets all followers of a user
func (r *UserRepository) GetFollowers(ctx context.Context, userID string) ([]model.User, error) {
	var followers []model.User
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Association("Followers").
		Find(&followers); err != nil {
		return nil, err
	}
	return followers, nil
}

// GetFollowing gets all users that a user is following
func (r *UserRepository) GetFollowing(ctx context.Context, userID string) ([]model.User, error) {
	var following []model.User
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Association("Following").
		Find(&following); err != nil {
		return nil, err
	}
	return following, nil
}

// IsFollowing checks if followerID is following followeeID
func (r *UserRepository) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN follows ON users.id = follows.followee_id").
		Where("users.id = ? AND follows.follower_id = ?", followeeID, followerID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetFollowerCount gets the number of followers
func (r *UserRepository) GetFollowerCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.followee_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetFollowingCount gets the number of users being followed
func (r *UserRepository) GetFollowingCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN follows ON users.id = follows.followee_id").
		Where("follows.follower_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SearchUsers searches for users by username or display name
func (r *UserRepository) SearchUsers(ctx context.Context, query string, limit, offset int) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).
		Where("username ILIKE ? OR display_name ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetAllUsers gets all users with pagination
func (r *UserRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// ExistsEmail checks if email exists
func (r *UserRepository) ExistsEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsUsername checks if username exists
func (r *UserRepository) ExistsUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserWithFollowersAndFollowing gets user with eager loaded followers and following
func (r *UserRepository) GetUserWithFollowersAndFollowing(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).
		Preload("Followers").
		Preload("Following").
		Where("id = ?", id).
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
