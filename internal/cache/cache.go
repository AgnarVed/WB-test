package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	// init map
	items := make(map[string]Item)
	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	// init GC when cleanup interval > 0
	if cleanupInterval > 0 {
		cache.StartGC()
	}
	return &cache
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64
	// if duration == 0, duration = default
	if duration == 0 {
		duration = c.defaultExpiration
	}
	// set expiration time
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	c.Lock()
	defer c.Unlock()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	item, found := c.items[key]
	// key not found
	if !found {
		return nil, false
	}
	// check exp time, if not - unexpired
	if item.Expiration > 0 {
		// if expired return nil
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Value, true
}

func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}
	delete(c.items, key)
	return nil
}

func (c *Cache) StartGC() {
	fmt.Println("Starting GC")
	go c.GC()
}

func (c *Cache) Exist(key string) (cond int) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.items[key]; !ok {
		return 0
	}
	return 1
}

func (c *Cache) GC() {
	for {
		// wait for time in cleanupInterval
		<-time.After(c.cleanupInterval)
		if c.items == nil {
			return
		}
		// check for expired items and delete them
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

// expiredKeys returns list of expired keys
func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()
	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}
	return
}

// clearItems deletes list of keys (expired)
func (c *Cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()
	for _, k := range keys {
		delete(c.items, k)
	}
}
