package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/logger"
	"marketplace/internal/models"
	"time"
)

type (
	UserStorage interface {
		CreateUser(ctx context.Context, login, password string) (*models.User, error)
		GetUserByLogin(ctx context.Context, login string) (*models.User, error)
		IsExists(ctx context.Context, login string) (bool, error)
	}

	userStorage struct {
		db *sql.DB
	}
)

func NewUserStorage(db *sql.DB) UserStorage {
	return &userStorage{db: db}
}

func (s *userStorage) CreateUser(ctx context.Context, login, password string) (*models.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO users (login, hash_password) VALUES ($1, $2) RETURNING id, created_at"
	var (
		id        int
		createdAt time.Time
	)

	err := s.db.QueryRowContext(ctxWithTimeout, query, login, password).Scan(&id, &createdAt)
	if err != nil {
		logger.Logger.Errorf("create user error: %v", err)
		return nil, err
	}

	return &models.User{
		ID:           id,
		Login:        login,
		HashPassword: password,
		CreatedAt:    createdAt,
	}, nil
}

func (s *userStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {

	isExists, err := s.IsExists(ctx, login)
	if err != nil {
		return nil, err
	}

	if !isExists {
		textErr := fmt.Sprintf("user with login - %s does not exists", login)
		return nil, errors.New(textErr)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.User
	query := "SELECT id, login, hash_password FROM users WHERE login = $1"

	err = s.db.QueryRowContext(ctxWithTimeout, query, login).Scan(&user.ID, &user.Login, &user.HashPassword)
	if err != nil {
		logger.Logger.Errorf("get user by email error: %v", err)
		return nil, err
	}

	return &user, nil
}

func (s *userStorage) IsExists(ctx context.Context, login string) (bool, error) {
	query := "SELECT EXISTS(SELECT * FROM users WHERE login = $1)"

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var isExists bool
	err := s.db.QueryRowContext(ctxWithTimeout, query, login).Scan(&isExists)
	if err != nil {
		logger.Logger.Errorf("exists check error: %v", err)
		return false, err
	}

	return isExists, nil
}
