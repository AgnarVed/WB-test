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

func (od order) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
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

func (od *order) CreateOrder(ctx context.Context, insert *models.Order) error {
	tx, err := od.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	err = od.orderDB.CreateOrder(ctx, tx, insert)
	if err != nil {
		od.RollbackTransaction(ctx, tx)
		return err
	}
	err = od.CommitTransaction(ctx, tx)
	if err != nil {
		od.RollbackTransaction(ctx, tx)
		return err
	}
	return nil
}

func (od *order) GetOrderList(ctx context.Context) ([]models.Order, error) {
	tx, err := od.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	orders, err := od.orderDB.GetOrderList(ctx, tx)

	if err != nil {
		return nil, err
	}
	return orders, nil
}

func NewOrderService(repos *repository.Repositories, cfg *config.Config) Order {
	return &order{
		orderDB:  repos.OrderDB,
		CommonDB: repos.CommonDB,
	}
}
