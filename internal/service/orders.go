package service

import (
	"context"
	"wbTest/internal/config"
	"wbTest/internal/models"
	"wbTest/internal/repository"
)

type order struct {
	repository.CommonDB
	orderDB repository.OrderDB
}

func (od order) GetOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
	tx, err := od.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	order, err := od.orderDB.GetOrderByID(ctx, tx, orderID)
	if err != nil {
		od.RollbackTransaction(ctx, tx)
		return nil, err
	}
	err = od.CommitTransaction(ctx, tx)
	if err != nil {
		od.RollbackTransaction(ctx, tx)
		return nil, err
	}
	return order, nil
}

func NewOrderService(repos *repository.Repositories, cfg *config.Config) Order {
	return &order{
		orderDB:  repos.OrderDB,
		CommonDB: repos.CommonDB,
	}
}
