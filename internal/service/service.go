package service

import (
	"context"
	"wbTest/internal/config"
	"wbTest/internal/models"
	"wbTest/internal/repository"
)

type Service struct {
	Order Order
	Cache Cache
}

type Order interface {
	GetOrderByID(ctx context.Context, orderID int) (*models.Order, error)
}

type Cache interface {
	GetOrderInCacheByID(ctx context.Context, orderID int) (*models.Order, error)
}

func NewService(repos *repository.Repositories, cfg *config.Config) *Service {
	return &Service{
		Order: NewOrderService(repos, cfg),
	}
}
