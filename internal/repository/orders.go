package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"wbTest/internal/models"
)

type order struct {
}

func (od *order) GetOrderByID(ctx context.Context, tx *sql.Tx, orderID string) (*models.Order, error) {
	query := `SELECT id, order_uid, order_info
	FROM orders
	WHERE id=$1;`

	row := tx.QueryRowContext(ctx, query, orderID)
	order := models.Order{}

	err := row.Scan(
		&order.ID,
		&order.OrderUID,
		&order.Data,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

type Msg struct {
	Message string `json:"msg"`
}

func convert(event interface{}) *models.Order {
	c := &models.Order{}
	mapstructure.Decode(event, &c)
	return c
}

// GetOrderList gets all orders from Database
func (od *order) GetOrderList(ctx context.Context, tx *sql.Tx) ([]models.Order, error) {
	query := `SELECT * FROM orders`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderList := make([]models.Order, 0)

	for rows.Next() {
		var tmpOrder models.Order
		err = rows.Scan(
			&tmpOrder.ID,
			tmpOrder.OrderUID,
			tmpOrder.Data,
		)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, tmpOrder)
	}
	return orderList, nil
}

func (od *order) CreateOrder(ctx context.Context, tx *sql.Tx, insert *models.OrderInput) error {
	query := `INSERT INTO orders (order_uid, order_info) VALUES ($1,$2)`
	var msg Msg
	err := json.Unmarshal(insert.OrderInfo, &msg)
	if err != nil {
		logrus.Fatal("Cannot Unmarshal Data", err)
	}
	dataJson, err := json.Marshal(msg)
	_, err = tx.ExecContext(ctx, query, insert.OrderUID, dataJson)
	if err != nil {
		return err
	}
	return nil
}

func NewOrderDB() OrderDB {
	return &order{}
}
