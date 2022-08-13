package service

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"wbTest/internal/models"
	"wbTest/internal/repository"
)

type cache struct {
	repository.CommonDB
	cached repository.CacheDB
}

func (c *cache) Get(ctx context.Context, key string) (*models.Order, error) {
	tx, err := c.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	order, err := c.cached.Get(ctx, tx, key)
	if err != nil {
		return nil, err
	}
	err = c.CommitTransaction(ctx, tx)
	if err != nil {
		c.RollbackTransaction(ctx, tx)
		return nil, err
	}
	return order, nil
}

func (c *cache) UploadCache(ctx context.Context) (bool, error) {
	tx, err := c.BeginTransaction(ctx)
	if err != nil {
		return false, err
	}
	err = c.cached.UploadCache(ctx, tx)
	if err != nil {
		c.RollbackTransaction(ctx, tx)
	}
	err = c.CommitTransaction(ctx, tx)
	if err != nil {
		c.RollbackTransaction(ctx, tx)
	}
	return true, nil
}

func (c *cache) Peek(key interface{}) (result *models.Order, ok bool) {
	order, ok := c.cached.Peek(key)
	return order, ok
}

func (c *cache) Set(key interface{}, value json.RawMessage) (ok bool) {
	ok = c.cached.Set(key, value)
	return ok
}

func (c *cache) NewCache(size int, store repository.OrderDB, logger *logrus.Logger) (*repository.Cache, error) {
	cache, err := c.cached.NewCache(size, store, logger)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewCacheService(repos *repository.Repositories) CacheService {
	return &cache{
		CommonDB: repos.CommonDB,
		cached:   repos.CacheDB,
	}
}
