package service

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"Backend-trainee-assignment-winter-2025/internal/repository"
	"golang.org/x/net/context"
)

type Authorization interface {
	SignUp(c context.Context, user *models.User) error
	SignIn(user *models.User, secret string, expiry int) (string, error)
	CheckIfUserExists(c context.Context, username string) error
}

type User interface {
	GetByUsername(username string) (*models.User, error)
}

type Service struct {
	Authorization
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthorizationService(repo),
	}
}
