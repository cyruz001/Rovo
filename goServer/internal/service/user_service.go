package service

import (
	"context"
	"errors"

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

func (s *UserService) Register(ctx context.Context, email, password string) (*model.User, error) {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    email,
		Password: hashed,
		Role:     "user",
	}

	return user, s.userRepo.Create(ctx, user)
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) FindByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}
