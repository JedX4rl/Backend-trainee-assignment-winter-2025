package repository

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"database/sql"
	"golang.org/x/net/context"
)

type User interface {
	SignUp(c context.Context, username, password string) (*models.User, error)
	GetUserByUsername(c context.Context, username string) (*models.User, error)
	GetInfo(c context.Context, userId int) (*models.InfoResponse, error)
}

type Transaction interface {
	SendMoney(c context.Context, senderId int, receiver string, amount int32) error
}
type Shop interface {
	BuyItem(c context.Context, userId int, item string) error
}

type Repository struct {
	User
	Transaction
	Shop
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:        NewUserRepository(db),
		Transaction: NewTransactionRepository(db),
		Shop:        NewShopRepository(db),
	}
}
