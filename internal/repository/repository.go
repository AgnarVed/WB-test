package repository

import (
	"context"
	"database/sql"
	"wbTest/internal/models"
	"wbTest/internal/repository/client"
)

type Repositories struct {
	OrderDB  OrderDB
	CommonDB CommonDB
}

type CommonDB interface {
	Get() *sql.DB
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
	CommitTransaction(ctx context.Context, tx *sql.Tx) error
	RollbackTransaction(ctx context.Context, tx *sql.Tx) error
}

type OrderDB interface {
	GetOrderByID(ctx context.Context, tx *sql.Tx, orderD int) (*models.Order, error)
}

func NewRepositories(psqlClient *client.PostgresClient) *Repositories {
	return &Repositories{
		OrderDB:  NewOrderDB(),
		CommonDB: NewCommonRepo(*psqlClient),
	}
}
