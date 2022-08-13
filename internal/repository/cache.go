package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"
	"wbTest/internal/models"
)

type Cache struct {
	cache   *lru.Cache
	size    int64
	storage OrderDB
	logger  *logrus.Logger
}

// NewCache constructs cache of the given size.
func (c *Cache) NewCache(size int, store OrderDB, logger *logrus.Logger) (*Cache, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	res := &Cache{
		cache:   cache,
		size:    int64(size),
		storage: store,
		logger:  logger,
	}
	return res, nil
}

// Set adds a value to the cache. Returns true if an eviction occurred.
func (c *Cache) Set(key interface{}, value json.RawMessage) (ok bool) {
	ok = c.cache.Add(key, value)
	return ok
}

// Get looks up a key's value from the cache. Updates item to "currently used"
func (c *Cache) Get(ctx context.Context, tx *sql.Tx, key string) (*models.Order, bool, error) {
	v, ok := c.cache.Get(key)
	if ok {
		c.logger.Info("got cache hit")
	}
	if !ok {
		o, err := c.storage.GetOrderByID(ctx, tx, key)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		if err != nil {
			return nil, false, err
		}
		c.Set(o.OrderUID, o.Data)
		v = o.Data
	}
	order := convert(v)
	return order, true, nil
}

// Peek looks up a key's value from the cache, but doesn't move element to the top of list.
func (c *Cache) Peek(key interface{}) (result *models.Order, ok bool) {
	order, ok := c.cache.Peek(key)
	result = convert(order)
	return result, ok
}

// UploadCache uploads all data from database to cache
func (c *Cache) UploadCache(ctx context.Context, tx *sql.Tx) error {
	orders, err := c.storage.GetOrderList(ctx, tx)
	if err != nil {
		return err
	}
	for _, v := range orders {
		c.Set(v.OrderUID, v.Data)
	}
	return nil
}
