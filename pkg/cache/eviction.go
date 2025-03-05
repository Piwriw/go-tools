package cache

import "container/list"

// EvictionStrategy 定义缓存淘汰策略接口
type EvictionStrategy interface {
	// Evict 淘汰数据
	Evict(c *Cache)
	// ShouldEvict 新的元素需要加入缓存时，是否需要淘汰旧数据
	ShouldEvict() bool
}

// NoEviction 无淘汰策略
type NoEviction struct{}

// Evict 不做任何淘汰
func (n *NoEviction) Evict(c *Cache) {}

// ShouldEvict 返回 false，表示不需要淘汰任何缓存
func (n *NoEviction) ShouldEvict() bool {
	return false
}

// LRU 淘汰策略
type LRU struct {
	cacheSize int
	cache     map[string]*list.Element
	evictList *list.List
}

// NewLRU 创建一个新的 LRU 策略对象
func NewLRU(cacheSize int) *LRU {
	return &LRU{
		cacheSize: cacheSize,
		cache:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

// Evict 执行 LRU 淘汰
func (l *LRU) Evict(c *Cache) {
	if len(c.data) >= l.cacheSize {
		// 淘汰最久未使用的元素
		elem := l.evictList.Back()
		if elem != nil {
			c.Delete(elem.Value.(string))
		}
	}
}

// ShouldEvict 判断是否需要淘汰数据
func (l *LRU) ShouldEvict() bool {
	return len(l.cache) >= l.cacheSize
}
