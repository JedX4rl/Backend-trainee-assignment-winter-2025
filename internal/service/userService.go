package service

import (
	customErrors "Backend-trainee-assignment-winter-2025/internal/errors"
	"Backend-trainee-assignment-winter-2025/internal/jwt"
	"Backend-trainee-assignment-winter-2025/internal/models"
	"Backend-trainee-assignment-winter-2025/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log/slog"
	"time"
)

type UserService struct {
	Repo   repository.User
	logger *slog.Logger
}

func (u *UserService) Auth(c context.Context, username, password string) (string, error) {
	u.logger.Debug("Auth func called")
	ctx, cancel := context.WithTimeout(c, time.Second*20)
	defer cancel()

	user, err := u.Repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			encryptedPassword, err := bcrypt.GenerateFromPassword(
				[]byte(password),
				bcrypt.DefaultCost,
			)
			if err != nil {
				u.logger.Error("Failed to generate password", "error", err)
				return "", customErrors.ErrCreateUser

			}
			if user, err = u.Repo.SignUp(ctx, username, string(encryptedPassword)); err != nil {
				u.logger.Error("Failed to sign up", "error", err)
				return "", customErrors.ErrCreateUser
			}
		} else {
			u.logger.Error("Unexpected error", "error", err)
			return "", customErrors.ErrCreateUser
		}
	}
	_, _ = []byte(user.Password), []byte(password)

	if valid := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); valid != nil {
		u.logger.Error("Invalid password", "error", valid)
		return "", customErrors.ErrInvalidPassword
	}
	token, err := accessToken.GenerateJWT(user, accessToken.TokenExpiry)
	if err != nil {
		u.logger.Error("Failed to generate JWT", "error", err)
		return "", fmt.Errorf("failed to generate token")
	}
	u.logger.Debug("Generated JWT finished successfully")
	return token, nil
}

func (u *UserService) GetUserByUsername(c context.Context, username string) (*models.User, error) {
	u.logger.Debug("GetUserByUsername called")
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	u.logger.Debug("GetUserByUsername finished successfully")
	return u.Repo.GetUserByUsername(ctx, username)
}

func (u *UserService) GetInfo(c context.Context, userId int) (*models.InfoResponse, error) {
	u.logger.Debug("GetUserInfo called")
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	u.logger.Debug("GetUserInfo finished successfully")
	return u.Repo.GetInfo(ctx, userId)
}

func NewUserService(repo repository.User, logger *slog.Logger) *UserService {
	return &UserService{
		Repo:   repo,
		logger: logger,
	}
}
