package service

import (
	accessToken "Backend-trainee-assignment-winter-2025/internal/jwt"
	"Backend-trainee-assignment-winter-2025/internal/models"
	"Backend-trainee-assignment-winter-2025/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

type AuthorizationService struct {
	Repo repository.Authorization
}

func (a *AuthorizationService) SignUp(c context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10) //TODO change timeout
	defer cancel()
	return a.Repo.SignUp(ctx, user)
}

func (a *AuthorizationService) SignIn(user *models.User, secret string, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &accessToken.Jwt{
		Username: user.Username,
		Id:       strconv.Itoa(user.Id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func (a *AuthorizationService) GetUserByUsername(c context.Context, username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return a.Repo.GetUserByUsername(ctx, username)
}

func (a *AuthorizationService) CheckIfUserExists(c context.Context, username string) error {
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	return a.Repo.CheckIfUserExists(ctx, username)
}

func NewAuthorizationService(repo repository.Authorization) *AuthorizationService {
	return &AuthorizationService{
		Repo: repo,
	}
}
