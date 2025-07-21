package services

import (
	"marketplace/internal/auth"
	"marketplace/internal/storage"
)

type Services struct {
	UserService UserService
	AdService   AdService
}

func NewServices(s *storage.Storage, jwt auth.JWTProvider) *Services {
	return &Services{
		UserService: NewUserService(s.User, jwt),
		AdService:   NewAdService(s.Ad),
	}
}
