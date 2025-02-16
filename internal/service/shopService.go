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

type ShopService struct {
	Repo   repository.Shop
	logger *slog.Logger
}

func (s ShopService) BuyItem(c context.Context, userId int, item string) error {
	s.logger.Debug("BuyItem called")
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	s.logger.Debug("BuyItem finished successfully")
	err := s.Repo.BuyItem(ctx, userId, item)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("not enough money")
		}
		return fmt.Errorf("internal error")
	}
	return nil
}

func NewShopService(repo repository.Shop, logger *slog.Logger) *ShopService {
	return &ShopService{
		Repo:   repo,
		logger: logger,
	}
}
