package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"goServer/internal/dto"
	"goServer/internal/model"
	"goServer/internal/repository"
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostService(pr repository.PostRepository, ur repository.UserRepository) *PostService {
	return &PostService{postRepo: pr, userRepo: ur}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(ctx context.Context, userID string, req dto.CreatePostReq) (*model.Post, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.Text == "" {
		return nil, errors.New("post text is required")
	}

	if len(req.Text) > 500 {
		return nil, errors.New("post text exceeds 500 characters")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	post := &model.Post{
		UserID:        userID,
		Text:          strings.TrimSpace(req.Text),
		CharCount:     len(req.Text),
		ReplyTo:       req.ReplyTo,
		IsQuote:       req.IsQuote,
		QuotedTweetID: req.QuotedPostID,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Add media if provided
	if len(req.MediaURLs) > 0 {
		for i, url := range req.MediaURLs {
			media := &model.Media{
				PostID:   post.ID,
				URL:      url,
				Position: i,
			}
			s.postRepo.AddMedia(ctx, media)
		}
	}

	return post, nil
}

// GetPostByID retrieves a post by ID
func (s *PostService) GetPostByID(ctx context.Context, postID string) (*model.Post, error) {
	if postID == "" {
		return nil, errors.New("post id is required")
	}

	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	return post, nil
}

// UpdatePost updates a post's text (only by creator)
func (s *PostService) UpdatePost(ctx context.Context, postID, userID, text string) (*model.Post, error) {
	if postID == "" || userID == "" {
		return nil, errors.New("post id and user id are required")
	}

	if text == "" || len(text) > 500 {
		return nil, errors.New("invalid post text")
	}

	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	if post.UserID != userID {
		return nil, errors.New("unauthorized: can only edit your own posts")
	}

	post.Text = strings.TrimSpace(text)
	post.CharCount = len(text)

	if err := s.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

// DeletePost deletes a post (only by creator or admin)
func (s *PostService) DeletePost(ctx context.Context, postID, userID string) error {
	if postID == "" || userID == "" {
		return errors.New("post id and user id are required")
	}

	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return errors.New("post not found")
	}

	if post.UserID != userID {
		return errors.New("unauthorized: can only delete your own posts")
	}

	if err := s.postRepo.Delete(ctx, postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// GetFeed retrieves the user's feed (posts from following)
func (s *PostService) GetFeed(ctx context.Context, userID string, limit, offset int) ([]model.Post, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.postRepo.GetFeed(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	return posts, nil
}

// GetUserTimeline retrieves all posts from a specific user
func (s *PostService) GetUserTimeline(ctx context.Context, username string, limit, offset int) ([]model.Post, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.postRepo.GetUserPosts(ctx, user.ID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user timeline: %w", err)
	}

	return posts, nil
}

// LikePost likes a post
func (s *PostService) LikePost(ctx context.Context, userID, postID string) error {
	if userID == "" || postID == "" {
		return errors.New("user id and post id are required")
	}

	// Verify post exists
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return errors.New("post not found")
	}

	// Check if already liked
	isLiked, err := s.postRepo.IsPostLiked(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check like status: %w", err)
	}
	if isLiked {
		return errors.New("already liked this post")
	}

	if err := s.postRepo.LikePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	return nil
}

// UnlikePost unlikes a post
func (s *PostService) UnlikePost(ctx context.Context, userID, postID string) error {
	if userID == "" || postID == "" {
		return errors.New("user id and post id are required")
	}

	// Check if already liked
	isLiked, err := s.postRepo.IsPostLiked(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check like status: %w", err)
	}
	if !isLiked {
		return errors.New("have not liked this post")
	}

	if err := s.postRepo.UnlikePost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	return nil
}

// RepostPost reposts a post
func (s *PostService) RepostPost(ctx context.Context, userID, postID string) error {
	if userID == "" || postID == "" {
		return errors.New("user id and post id are required")
	}

	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return errors.New("post not found")
	}

	if post.UserID == userID {
		return errors.New("cannot repost your own post")
	}

	isReposted, err := s.postRepo.IsPostReposted(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check repost status: %w", err)
	}
	if isReposted {
		return errors.New("already reposted this post")
	}

	if err := s.postRepo.RepostPost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to repost: %w", err)
	}

	return nil
}

// UndoRepost removes a repost
func (s *PostService) UndoRepost(ctx context.Context, userID, postID string) error {
	if userID == "" || postID == "" {
		return errors.New("user id and post id are required")
	}

	isReposted, err := s.postRepo.IsPostReposted(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check repost status: %w", err)
	}
	if !isReposted {
		return errors.New("have not reposted this post")
	}

	if err := s.postRepo.UndoRepost(ctx, userID, postID); err != nil {
		return fmt.Errorf("failed to undo repost: %w", err)
	}

	return nil
}

// IsPostLiked checks if a post is liked by a user
func (s *PostService) IsPostLiked(ctx context.Context, userID, postID string) (bool, error) {
	return s.postRepo.IsPostLiked(ctx, userID, postID)
}

// IsPostReposted checks if a post is reposted by a user
func (s *PostService) IsPostReposted(ctx context.Context, userID, postID string) (bool, error) {
	return s.postRepo.IsPostReposted(ctx, userID, postID)
}

// GetPostLikes gets all users who liked a post
func (s *PostService) GetPostLikes(ctx context.Context, postID string, limit, offset int) ([]model.User, error) {
	if postID == "" {
		return nil, errors.New("post id is required")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.postRepo.GetPostLikes(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get post likes: %w", err)
	}

	return users, nil
}

// GetPostReposts gets all users who reposted a post
func (s *PostService) GetPostReposts(ctx context.Context, postID string, limit, offset int) ([]model.User, error) {
	if postID == "" {
		return nil, errors.New("post id is required")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.postRepo.GetPostReposts(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get post reposts: %w", err)
	}

	return users, nil
}

// GetReplies gets all replies to a post
func (s *PostService) GetReplies(ctx context.Context, postID string, limit, offset int) ([]model.Post, error) {
	if postID == "" {
		return nil, errors.New("post id is required")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.postRepo.GetReplies(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}

	return posts, nil
}

// SearchPosts searches for posts by text
func (s *PostService) SearchPosts(ctx context.Context, query string, limit, offset int) ([]model.Post, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.postRepo.SearchPosts(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}

	return posts, nil
}

// DeletePostAdmin deletes a post (admin only)
func (s *PostService) DeletePostAdmin(ctx context.Context, postID string) error {
	if postID == "" {
		return errors.New("post id is required")
	}

	if err := s.postRepo.Delete(ctx, postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
