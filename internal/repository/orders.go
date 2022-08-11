package repository

import (
	"context"
	"database/sql"
	"wbTest/internal/models"
)

type order struct {
}

func (od *order) GetOrderByID(ctx context.Context, tx *sql.Tx, orderID int) (*models.Order, error) {
	query := `SELECT id, order_uid, data
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

func (od *order) CreateOrder(ctx context.Context, tx *sql.Tx, insert *models.Order) error {
	query := `INSERT INTO orders VALUES ($1,$2,$3)`
	_, err := tx.ExecContext(ctx, query, insert.ID, insert.OrderUID, insert.Data)
	if err != nil {
		return err
	}

	return nil
}

func NewOrderDB() OrderDB {
	return &order{}
}
