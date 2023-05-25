package repository

import (
	"context"
	"database/sql"
	"wbTest/internal/repository/client"
)

type commonDBRepo struct {
	client.PostgresClient
}

func NewCommonRepo(postgresClient client.PostgresClient) CommonDB {
	return &commonDBRepo{
		PostgresClient: postgresClient,
	}
}

func (c *commonDBRepo) Get() *sql.DB {
	cl, _ := c.GetClient()
	return cl
}

func (c *commonDBRepo) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	cl, err := c.GetClient()
	if err != nil {
		return nil, err
	}

	opts := sql.TxOptions{
		Isolation: sql.LevelDefault,
	}

	tx, err := cl.BeginTx(ctx, &opts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *commonDBRepo) CommitTransaction(_ context.Context, tx *sql.Tx) error {
	return tx.Commit()
}

func (c *commonDBRepo) RollbackTransaction(_ context.Context, tx *sql.Tx) error {
	return tx.Rollback()
}
