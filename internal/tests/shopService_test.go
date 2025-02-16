package tests

import (
	"Backend-trainee-assignment-winter-2025/internal/service"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"log/slog"
	"os"
	"testing"
)

type MockShopRepository struct {
	mock.Mock
	db *sql.DB
}

func (m *MockShopRepository) BuyItem(c context.Context, userId int, item string) error {
	args := m.Called(c, userId, item)
	return args.Error(0)
}

func TestShopService_BuyItem_NotEnoughMoney(t *testing.T) {
	mockRepo := new(MockShopRepository)
	mockLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) // Мок логгера
	shopService := service.NewShopService(mockRepo, mockLogger)

	mockRepo.On("BuyItem", mock.Anything, 1, "T-shirt").Return(sql.ErrNoRows)

	ctx := context.Background()

	err := shopService.BuyItem(ctx, 1, "T-shirt")

	assert.Error(t, err)
	assert.Equal(t, "not enough money", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestShopService_BuyItem_InternalError(t *testing.T) {
	mockRepo := new(MockShopRepository)
	mockLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	shopService := service.NewShopService(mockRepo, mockLogger)

	mockRepo.On("BuyItem", mock.Anything, 1, "T-shirt").Return(fmt.Errorf("internal error"))

	ctx := context.Background()

	err := shopService.BuyItem(ctx, 1, "T-shirt")

	assert.Error(t, err)
	assert.Equal(t, "internal error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestShopService_BuyItem_Success(t *testing.T) {
	mockRepo := new(MockShopRepository)
	mockLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	shopService := service.NewShopService(mockRepo, mockLogger)

	mockRepo.On("BuyItem", mock.Anything, 1, "T-shirt").Return(nil)

	ctx := context.Background()

	err := shopService.BuyItem(ctx, 1, "T-shirt")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
