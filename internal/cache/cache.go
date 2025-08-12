package cache

import (
	"my-shop/internal/model"
	"sync"
)

type OrderCache struct {
	mu    sync.RWMutex
	items map[string]*model.Order
}

func New() *OrderCache {
	return &OrderCache{
		items: make(map[string]*model.Order),
	}
}

func (c *OrderCache) Get(id string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.items[id]
	return order, ok
}

func (c *OrderCache) Set(id string, order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[id] = order
}

func (c *OrderCache) All() map[string]*model.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	copy := make(map[string]*model.Order, len(c.items))
	for k, v := range c.items {
		copy[k] = v
	}
	return copy
}
