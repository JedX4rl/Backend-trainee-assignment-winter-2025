package repository

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"database/sql"
	"golang.org/x/net/context"
)

type Authorization interface {
	SignUp(c context.Context, user *models.User) error
	CheckIfUserExists(c context.Context, username string) error
}

type User interface {
	GetByUsername(c context.Context, username string) (*models.User, error)
}

type Repository struct {
	Authorization
	User
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthorizationRepository(db),
	}
}
