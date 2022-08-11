package client

import (
	"database/sql"
	"errors"
)

type PostgresClient interface {
	GetClient() (*sql.DB, error)
}

type postgresClient struct {
	db *sql.DB
}

func (p postgresClient) GetClient() (*sql.DB, error) {
	if p.db == nil {
		return nil, errors.New("pg is nil")
	}
	return p.db, nil
}

func NewPostgresClient(db *sql.DB) PostgresClient {
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(0)
	return &postgresClient{
		db: db,
	}
}
