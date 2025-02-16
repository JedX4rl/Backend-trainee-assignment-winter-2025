package service

import (
	"Backend-trainee-assignment-winter-2025/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"log/slog"
	"time"
)

type TransactionService struct {
	Repo   repository.Transaction
	logger *slog.Logger
}

func (t TransactionService) SendMoney(c context.Context, senderId int, receiver string, amount int32) error {
	t.logger.Debug("SendMoney called")
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()

	err := t.Repo.SendMoney(ctx, senderId, receiver, amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("not enough money or invalid user")
		}
		return fmt.Errorf("internal error")
	}

	t.logger.Debug("SendMoney finished successfully")
	return nil
}

func NewTransactionService(repo repository.Transaction, logger *slog.Logger) *TransactionService {
	return &TransactionService{
		Repo:   repo,
		logger: logger,
	}
}
