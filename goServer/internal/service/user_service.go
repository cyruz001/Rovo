package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"goServer/internal/dto"
	"goServer/internal/model"
	"goServer/internal/repository"
	"goServer/pkg/utils"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{userRepo: r}
}

// Register creates a new user account
func (s *UserService) Register(ctx context.Context, email, username, password string) (*model.User, error) {
	// Validate input
	if email == "" || username == "" || password == "" {
		return nil, errors.New("email, username, and password are required")
	}

	// Check if email already exists
	existsEmail, err := s.userRepo.ExistsEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existsEmail {
		return nil, errors.New("email already registered")
	}

	// Check if username already exists
	existsUsername, err := s.userRepo.ExistsUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existsUsername {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
		Role:     "USER",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Authenticate verifies user credentials
// email and password
// func (s *UserService) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
// 	if email == "" || password == "" {
// 		return nil, errors.New("email and password are required")
// 	}

// 	user, err := s.userRepo.FindByEmail(ctx, email)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to find user: %w", err)
// 	}
// 	if user == nil {
// 		return nil, errors.New("user not found")
// 	}

// 	if err := utils.CheckPassword(password, user.Password); err != nil {
// 		return nil, errors.New("invalid credentials")
// 	}

//		return user, nil
//	}
//
// username and password
func (s *UserService) Authenticate(ctx context.Context, username, password string) (*model.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
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

	return user, nil
}

// UpdateUser updates user profile information
func (s *UserService) UpdateUser(ctx context.Context, userID string, req dto.UserUpdateReq) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update only non-empty fields
	if req.DisplayName != "" {
		user.DisplayName = strings.TrimSpace(req.DisplayName)
	}
	if req.Bio != "" {
		user.Bio = strings.TrimSpace(req.Bio)
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user account
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user id is required")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// FollowUser creates a follow relationship
func (s *UserService) FollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == "" || followeeID == "" {
		return errors.New("follower id and followee id are required")
	}

	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}

	// Check if followee exists
	followee, err := s.userRepo.FindByID(ctx, followeeID)
	if err != nil {
		return fmt.Errorf("failed to find followee: %w", err)
	}
	if followee == nil {
		return errors.New("followee not found")
	}

	// Check if already following
	isFollowing, err := s.userRepo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}
	if isFollowing {
		return errors.New("already following this user")
	}

	if err := s.userRepo.FollowUser(ctx, followerID, followeeID); err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}

	return nil
}

// UnfollowUser removes a follow relationship
func (s *UserService) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == "" || followeeID == "" {
		return errors.New("follower id and followee id are required")
	}

	// Check if following
	isFollowing, err := s.userRepo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}
	if !isFollowing {
		return errors.New("not following this user")
	}

	if err := s.userRepo.UnfollowUser(ctx, followerID, followeeID); err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}

	return nil
}

// GetFollowers retrieves all followers of a user
func (s *UserService) GetFollowers(ctx context.Context, username string) ([]model.User, error) {
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

	followers, err := s.userRepo.GetFollowers(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	return followers, nil
}

// GetFollowing retrieves all users that a user is following
func (s *UserService) GetFollowing(ctx context.Context, username string) ([]model.User, error) {
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

	following, err := s.userRepo.GetFollowing(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	return following, nil
}

// GetFollowerCount gets the number of followers
func (s *UserService) GetFollowerCount(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, errors.New("user id is required")
	}

	count, err := s.userRepo.GetFollowerCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get follower count: %w", err)
	}

	return count, nil
}

// GetFollowingCount gets the number of users being followed
func (s *UserService) GetFollowingCount(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, errors.New("user id is required")
	}

	count, err := s.userRepo.GetFollowingCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get following count: %w", err)
	}

	return count, nil
}

// SearchUsers searches for users by username or display name
func (s *UserService) SearchUsers(ctx context.Context, query string, limit, offset int) ([]model.User, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.SearchUsers(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.GetAllUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

// UpdateUserRole updates a user's role (admin only)
func (s *UserService) UpdateUserRole(ctx context.Context, userID, role string) (*model.User, error) {
	if userID == "" || role == "" {
		return nil, errors.New("user id and role are required")
	}

	// Validate role
	if role != "USER" && role != "ADMIN" {
		return nil, errors.New("invalid role")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Role = role

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	return user, nil
}

// IsFollowing checks if one user follows another
func (s *UserService) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	if followerID == "" || followeeID == "" {
		return false, errors.New("follower id and followee id are required")
	}

	return s.userRepo.IsFollowing(ctx, followerID, followeeID)
}
