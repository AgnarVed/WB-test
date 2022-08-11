package repository

import (
	"context"
	"database/sql"
	"wbTest/internal/models"
)

type order struct {
}

func NewOrderDB() OrderDB {
	return &order{}
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
