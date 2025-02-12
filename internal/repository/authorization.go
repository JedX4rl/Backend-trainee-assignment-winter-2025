package repository

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"database/sql"
	"errors"
	"golang.org/x/net/context"
)

type AuthorizationRepository struct {
	db *sql.DB
}

func (a AuthorizationRepository) SignUp(c context.Context, user *models.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	row := a.db.QueryRowContext(c, query, user.Username, user.Password)
	return row.Err()
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

func (a AuthorizationRepository) GetUserByUsername(c context.Context, username string) (*models.User, error) {

	query := `SELECT * FROM users WHERE username = $1`
	row := a.db.QueryRowContext(c, query, username)

	if err := row.Err(); err != nil {
		return nil, err
	}

	var user models.User

	if err := row.Scan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func NewAuthorizationRepository(db *sql.DB) *AuthorizationRepository {
	return &AuthorizationRepository{
		db: db,
	}
}
