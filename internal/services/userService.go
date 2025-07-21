package services

import (
	"context"
	"errors"
	"fmt"
	"marketplace/internal/auth"
	"marketplace/internal/logger"
	"marketplace/internal/models"
	"marketplace/internal/storage"
)

type (
	UserService interface {
		Register(ctx context.Context, login, password string) (*models.User, error)
		Login(ctx context.Context, login, password string) (string, error)
	}

	userService struct {
		storage  storage.UserStorage
		jwtUtils auth.JWTProvider
	}
)

func NewUserService(storage storage.UserStorage, jwtUtils auth.JWTProvider) UserService {
	return &userService{
		storage:  storage,
		jwtUtils: jwtUtils,
	}
}

func (s *userService) Register(ctx context.Context, login, password string) (*models.User, error) {

	exists, err := s.storage.IsExists(ctx, login)
	if err != nil {
		logger.Logger.Errorf("%v", err)
		return nil, err
	}
	if exists {
		textErr := fmt.Sprintf("user with login: \"%s\" alredy exists", login)
		return nil, errors.New(textErr)
	}

	HashPassword, err := HashPassword(password)
	if err != nil {
		logger.Logger.Errorf("password hashing error: %v", err)
		return nil, fmt.Errorf("password hashing error: %v", err)
	}

	return s.storage.CreateUser(ctx, login, HashPassword)
}

func (s *userService) Login(ctx context.Context, login, password string) (string, error) {

	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		logger.Logger.Errorf("invalid login or password: %v", err)
		return "", err
	}

	err = ValidatePassword(user.HashPassword, password)
	if err != nil {
		logger.Logger.Errorf("invalid login or password: %v", err)
		return "", err
	}

	token, err := s.jwtUtils.GenerateToken(user.ID)
	if err != nil {
		logger.Logger.Errorf("jwt token creating error: %v", err)
	}

	return token, nil
}
