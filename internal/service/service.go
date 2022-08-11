package service

import (
	"context"
	"wbTest/internal/config"
	"wbTest/internal/models"
	"wbTest/internal/repository"
)

type Service struct {
	Order Order
}

type Order interface {
	GetOrderByID(ctx context.Context, orderID int) (*models.Order, error)
	CreateOrder(ctx context.Context, insert *models.Order) error
}

func NewService(repos *repository.Repositories, cfg *config.Config) *Service {
	return &Service{
		Order: NewOrderService(repos, cfg),
	}
}
