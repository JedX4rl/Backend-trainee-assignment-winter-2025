package service

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"Backend-trainee-assignment-winter-2025/internal/repository"
	"golang.org/x/net/context"
	"log/slog"
)

type User interface {
	Auth(c context.Context, username, password string) (models.AuthResponse, error)
	GetUserByUsername(c context.Context, username string) (*models.User, error)
	GetInfo(c context.Context, userId int) (*models.InfoResponse, error)
}

type Transaction interface {
	SendMoney(c context.Context, senderId int, receiver string, amount int32) error
}

type Shop interface {
	BuyItem(c context.Context, userId int, item string) error
}

type Service struct {
	User
	Transaction
	Shop
}

func NewService(repo *repository.Repository, logger *slog.Logger) *Service {
	return &Service{
		User:        NewUserService(repo, logger),
		Transaction: NewTransactionService(repo, logger),
		Shop:        NewShopService(repo, logger),
	}
}
