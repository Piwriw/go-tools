package cache

import (
	"sync"
)

// Cache 缓存结构
type Cache struct {
	data map[string]any
	// 缓存的最大容量
	capacity int
	mu       sync.RWMutex
	// 采用的淘汰策略
	evictionStrategy EvictionStrategy
}

// NewCache 创建缓存实例
func NewCache(capacity int, evictionStrategy EvictionStrategy) *Cache {
	if evictionStrategy == nil {
		evictionStrategy = &NoEviction{}
	}
	return &Cache{
		data:             make(map[string]any),
		evictionStrategy: evictionStrategy,
		capacity:         capacity,
	}
}

// Update 更新缓存中的数据，仅当值不同时才更新
func (c *Cache) Update(key string, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 获取缓存中现有的值
	existingValue, exists := c.data[key]
	if !exists {
		return ErrKeyNotFound
	}

	// 如果现有值和新值不一样，才更新
	if existingValue != value {
		c.data[key] = value
	}
	return nil
}

// Set 存储数据到缓存
func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果缓存满了，并且需要进行淘汰操作
	if len(c.data) >= c.capacity {
		// 如果配置了淘汰策略，则进行淘汰
		if c.evictionStrategy.ShouldEvict() {
			c.evictionStrategy.Evict(c)
		}
	}

	c.data[key] = value
}

// Get 从缓存中获取数据
func (c *Cache) Get(key string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.data[key]
	if !exists {
		return nil, ErrKeyNotFound
	}
	return value, nil
}

// Delete 删除缓存中的数据
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear 清空缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]any)
}
