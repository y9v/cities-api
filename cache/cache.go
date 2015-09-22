package cache

import (
	"sync"
)

type Cache struct {
	lock sync.RWMutex
	data map[string]interface{}
}

func New() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	d, ok := c.data[key]
	return d, ok
}

func (c *Cache) Set(key string, d interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	c.data[key] = d
}
