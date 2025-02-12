package repository

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"database/sql"
	"golang.org/x/net/context"
)

type AuthorizationRepository struct {
	db *sql.DB
}

func (a AuthorizationRepository) SignUp(c context.Context, user *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (a AuthorizationRepository) CheckIfUserExists(c context.Context, username string) error {
	//TODO should I add a transaction?
	query := `SELECT username FROM users WHERE username = $1`
	row := a.db.QueryRow(query, username)
	if err := row.Err(); err != nil {
		return err //TODO think about errors
	}
	return nil
}

func NewAuthorizationRepository(db *sql.DB) *AuthorizationRepository {
	return &AuthorizationRepository{
		db: db,
	}
}
