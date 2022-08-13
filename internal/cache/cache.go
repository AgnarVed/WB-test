package cache

//// LRUCache implements a non-thread safe fixed size LRU cache
//type LRUCache struct {
//	size      int
//	evictList *list.List
//	onEvict   EvictCallback
//	//sync.RWMutex
//	//defaultExpiration time.Duration
//	//cleanupInterval   time.Duration
//	items map[interface{}]*list.Element
//}
//
//// item is used to hold a value in the evictList
//type item struct {
//	key   interface{}
//	value interface{}
//}
//
//// EvictCallback is used to get a callback when a cache entry is evicted.
//type EvictCallback func(key interface{}, value interface{})
//
//// NewCache constructs cache of the given size.
//func NewCache(size int, onEvict EvictCallback) (*LRUCache, error) {
//	// init map
//	if size <= 0 {
//		return nil, errors.New("size must provide a positive number")
//	}
//	cache := &LRUCache{
//		size:      size,
//		evictList: list.New(),
//		onEvict:   onEvict,
//		//RWMutex:   sync.RWMutex{},
//		items: make(map[interface{}]*list.Element),
//	}
//	return cache, nil
//}
//
//// Clear is used to completely clear the cache.
//func (c *LRUCache) Clear() {
//	for k, v := range c.items {
//		if c.onEvict != nil {
//			c.onEvict(k, v.Value.(*item).value)
//		}
//		delete(c.items, k)
//	}
//	c.evictList.Init()
//}
//
//// Set adds a value to the cache. Returns true if an eviction occurred.
//func (c *LRUCache) Set(key, value interface{}) (evicted bool) {
//	// Check for existing item
//	if ent, ok := c.items[key]; ok {
//		c.evictList.MoveToFront(ent)
//		ent.Value.(*item).value = value
//		return false
//	}
//
//	// Add new item
//	ent := &item{key, value}
//	entry := c.evictList.PushFront(ent)
//	c.items[key] = entry
//
//	evict := c.evictList.Len() > c.size
//	// Verify size not exceeded
//	if evict {
//		c.deleteOldest()
//	}
//	return evict
//}
//
//// Get looks up a key's value from the cache, and moves element to the front of list.
//func (c *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
//	if ent, ok := c.items[key]; ok {
//		c.evictList.MoveToFront(ent)
//		if ent.Value.(*item) == nil {
//			return nil, false
//		}
//		return ent.Value.(*item).value, true
//	}
//	return
//}
//
//// Peek looks up a key's value from the cache, but doesn't move element to the top of list.
//func (c *LRUCache) Peek(key interface{}) (value interface{}, ok bool) {
//	var ent *list.Element
//	if ent, ok = c.items[key]; ok {
//		return ent.Value.(*item).value, true
//	}
//	return nil, ok
//}
//
//// Delete removes the provided key from cache,returning if the key was contained.
//func (c *LRUCache) Delete(key interface{}) bool {
//	//c.Lock()
//	//defer c.Unlock()
//	if ent, ok := c.items[key]; ok {
//		c.deleteElement(ent)
//		return true
//	}
//	return false
//}
//
//// DeleteOldest removes the oldest item from the cache.
//func (c *LRUCache) DeleteOldest() (key, value interface{}, ok bool) {
//	ent := c.evictList.Back()
//	if ent != nil {
//		c.deleteElement(ent)
//		kv := ent.Value.(*item)
//		return kv.key, kv.value, true
//	}
//	return nil, nil, false
//}
//
//// Exist checks if item is in cache.
//func (c *LRUCache) Exist(key string) (ok bool) {
//	//c.Lock()
//	//defer c.Unlock()
//	_, ok = c.items[key]
//	return ok
//}
//
//// GetOldest returns the oldest item.
//func (c *LRUCache) GetOldest() (key, value interface{}, ok bool) {
//	ent := c.evictList.Back()
//	if ent != nil {
//		kv := ent.Value.(*item)
//		return kv.key, kv.value, true
//	}
//	return nil, nil, false
//}
//
//// Len returns the number of items in the cache.
//func (c *LRUCache) Len() int {
//	return c.evictList.Len()
//}
//
//// Keys returns a slice of the keys in the cache, from oldest to newest.
//func (c *LRUCache) Keys() []interface{} {
//	keys := make([]interface{}, len(c.items))
//	i := 0
//	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
//		keys[i] = ent.Value.(*item).key
//		i++
//	}
//	return keys
//}
//
//// Resize changes the cache size.
//func (c *LRUCache) Resize(size int) (evicted int) {
//	diff := c.Len() - size
//	if diff < 0 {
//		diff = 0
//	}
//	for i := 0; i < diff; i++ {
//		c.deleteOldest()
//	}
//	c.size = size
//	return diff
//}
//
//// deleteElement is used to remove a given list element from the cache.
//func (c *LRUCache) deleteElement(e *list.Element) {
//	c.evictList.Remove(e)
//	kv := e.Value.(*item)
//	delete(c.items, kv.key)
//	if c.onEvict != nil {
//		c.onEvict(kv.key, kv.value)
//	}
//}
//
//// deleteOldest removes the oldest item from the cache.
//func (c *LRUCache) deleteOldest() {
//	ent := c.evictList.Back()
//	if ent != nil {
//		c.deleteElement(ent)
//	}
//}
