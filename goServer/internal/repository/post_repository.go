package repository

import (
	"context"
	"errors"

	"goServer/internal/model"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create creates a new post
func (r *PostRepository) Create(ctx context.Context, p *model.Post) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// FindByID finds a post by ID
func (r *PostRepository) FindByID(ctx context.Context, id string) (*model.Post, error) {
	var p model.Post
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Media").
		Where("id = ?", id).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// Update updates a post
func (r *PostRepository) Update(ctx context.Context, p *model.Post) error {
	return r.db.WithContext(ctx).Model(p).Updates(p).Error
}

// Delete deletes a post
func (r *PostRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Post{}).Error
}

// GetUserPosts gets all posts from a user
func (r *PostRepository) GetUserPosts(ctx context.Context, userID string, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Media").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// GetFeed gets posts from users that the current user follows
func (r *PostRepository) GetFeed(ctx context.Context, userID string, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Media").
		Where("user_id IN (SELECT followee_id FROM follows WHERE follower_id = ?)", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// AddMedia adds media to a post
func (r *PostRepository) AddMedia(ctx context.Context, m *model.Media) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// LikePost adds a like to a post
func (r *PostRepository) LikePost(ctx context.Context, userID, postID string) error {
	like := &model.Like{
		UserID: userID,
		PostID: postID,
	}
	return r.db.WithContext(ctx).Create(like).Error
}

// UnlikePost removes a like from a post
func (r *PostRepository) UnlikePost(ctx context.Context, userID, postID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&model.Like{}).Error
}

// IsPostLiked checks if a post is liked by a user
func (r *PostRepository) IsPostLiked(ctx context.Context, userID, postID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Like{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetLikeCount gets the number of likes on a post
func (r *PostRepository) GetLikeCount(ctx context.Context, postID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Like{}).
		Where("post_id = ?", postID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetPostLikes gets all users who liked a post
func (r *PostRepository) GetPostLikes(ctx context.Context, postID string, limit, offset int) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).
		Joins("JOIN likes ON users.id = likes.user_id").
		Where("likes.post_id = ?", postID).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// RepostPost reposts a post
func (r *PostRepository) RepostPost(ctx context.Context, userID, postID string) error {
	repost := &model.Repost{
		UserID: userID,
		PostID: postID,
	}
	return r.db.WithContext(ctx).Create(repost).Error
}

// UndoRepost removes a repost
func (r *PostRepository) UndoRepost(ctx context.Context, userID, postID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&model.Repost{}).Error
}

// IsPostReposted checks if a post is reposted by a user
func (r *PostRepository) IsPostReposted(ctx context.Context, userID, postID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Repost{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetRepostCount gets the number of reposts on a post
func (r *PostRepository) GetRepostCount(ctx context.Context, postID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Repost{}).
		Where("post_id = ?", postID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetPostReposts gets all users who reposted a post
func (r *PostRepository) GetPostReposts(ctx context.Context, postID string, limit, offset int) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).
		Joins("JOIN reposts ON users.id = reposts.user_id").
		Where("reposts.post_id = ?", postID).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetReplies gets all replies to a post
func (r *PostRepository) GetReplies(ctx context.Context, postID string, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Media").
		Where("reply_to = ?", postID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// GetReplyCount gets the number of replies to a post
func (r *PostRepository) GetReplyCount(ctx context.Context, postID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("reply_to = ?", postID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SearchPosts searches for posts by text
func (r *PostRepository) SearchPosts(ctx context.Context, query string, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Media").
		Where("text ILIKE ?", "%"+query+"%").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// GetAllPosts gets all posts with pagination
func (r *PostRepository) GetAllPosts(ctx context.Context, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Media").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
