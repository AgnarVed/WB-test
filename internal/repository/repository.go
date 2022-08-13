package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"wbTest/internal/models"
	"wbTest/internal/repository/client"
)

type Repositories struct {
	OrderDB  OrderDB
	CommonDB CommonDB
	CacheDB  CacheDB
}

type CommonDB interface {
	Get() *sql.DB
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
	CommitTransaction(ctx context.Context, tx *sql.Tx) error
	RollbackTransaction(ctx context.Context, tx *sql.Tx) error
}

type OrderDB interface {
	GetOrderByID(ctx context.Context, tx *sql.Tx, orderID string) (*models.Order, error)
	CreateOrder(ctx context.Context, tx *sql.Tx, insert *models.Order) error
	// GetOrderList gets all orders from Database
	GetOrderList(ctx context.Context, tx *sql.Tx) ([]models.Order, error)
}

type CacheDB interface {
	NewCache(size int, store OrderDB, logger *logrus.Logger) (*Cache, error)
	Get(ctx context.Context, tx *sql.Tx, orderID string) (*models.Order, error)
	Set(key interface{}, value json.RawMessage) (ok bool)
	UploadCache(ctx context.Context, tx *sql.Tx) error
	Peek(key interface{}) (result *models.Order, ok bool)
}

func NewRepositories(psqlClient *client.PostgresClient) *Repositories {
	return &Repositories{
		OrderDB:  NewOrderDB(),
		CommonDB: NewCommonRepo(*psqlClient),
	}
}
