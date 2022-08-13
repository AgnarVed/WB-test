package service

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"wbTest/internal/config"
	"wbTest/internal/models"
	"wbTest/internal/repository"
)

type Service struct {
	Order Order
	Cache CacheService
}

type Order interface {
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	CreateOrder(ctx context.Context, insert *models.Order) error
	GetOrderList(ctx context.Context) ([]models.Order, error)
}

type CacheService interface {
	NewCache(size int, store repository.OrderDB, logger *logrus.Logger) (*repository.Cache, error)
	Get(ctx context.Context, key string) (*models.Order, error)
	Set(key interface{}, value json.RawMessage) (ok bool)
	UploadCache(ctx context.Context) (bool, error)
	Peek(key interface{}) (result *models.Order, ok bool)
}

func NewService(repos *repository.Repositories, cfg *config.Config) *Service {
	return &Service{
		Order: NewOrderService(repos, cfg),
		Cache: NewCacheService(repos),
	}
}
