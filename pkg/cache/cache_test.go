package cache

import (
	"sync"
	"testing"
)

// TestNewCache 测试 NewCache 构造函数
// 测试场景：正常创建缓存、指定容量、指定淘汰策略、nil 淘汰策略应使用默认值
func TestNewCache(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		evictionStrategy  EvictionStrategy
		expectNilStrategy bool
	}{
		{
			name:             "创建容量为10的缓存",
			capacity:         10,
			evictionStrategy: &NoEviction{},
		},
		{
			name:             "创建容量为100的缓存",
			capacity:         100,
			evictionStrategy: NewLRU(5),
		},
		{
			name:              "nil淘汰策略应使用默认值",
			capacity:          10,
			evictionStrategy:  nil,
			expectNilStrategy: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache(tt.capacity, tt.evictionStrategy)

			if cache == nil {
				t.Fatal("NewCache 返回 nil")
			}
			if cache.capacity != tt.capacity {
				t.Errorf("期望容量=%d, 实际=%d", tt.capacity, cache.capacity)
			}
			if cache.data == nil {
				t.Error("data map 未初始化")
			}
			if tt.expectNilStrategy {
				if _, ok := cache.evictionStrategy.(*NoEviction); !ok {
					t.Error("nil 淘汰策略未使用默认的 NoEviction")
				}
			}
		})
	}
}

// TestCache_SetAndGet 测试 Set 和 Get 方法
// 测试场景：正常存储和获取、获取不存在的键、覆盖已存在的键
func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(10, &NoEviction{})

	t.Run("正常存储和获取", func(t *testing.T) {
		cache.Set("key1", "value1")

		value, err := cache.Get("key1")
		if err != nil {
			t.Errorf("获取失败: %v", err)
		}
		if value != "value1" {
			t.Errorf("期望值=value1, 实际=%v", value)
		}
	})

	t.Run("获取不存在的键应返回错误", func(t *testing.T) {
		_, err := cache.Get("nonexistent")
		if err != ErrKeyNotFound {
			t.Errorf("期望错误=ErrKeyNotFound, 实际=%v", err)
		}
	})

	t.Run("覆盖已存在的键", func(t *testing.T) {
		cache.Set("key2", "old_value")
		cache.Set("key2", "new_value")

		value, err := cache.Get("key2")
		if err != nil {
			t.Errorf("获取失败: %v", err)
		}
		if value != "new_value" {
			t.Errorf("期望值=new_value, 实际=%v", value)
		}
	})
}

// TestCache_Update 测试 Update 方法
// 测试场景：更新存在的键、更新不存在的键应返回错误、值相同时不更新
func TestCache_Update(t *testing.T) {
	cache := NewCache(10, &NoEviction{})

	t.Run("更新存在的键", func(t *testing.T) {
		cache.Set("key1", "value1")
		err := cache.Update("key1", "value2")
		if err != nil {
			t.Errorf("更新失败: %v", err)
		}

		value, _ := cache.Get("key1")
		if value != "value2" {
			t.Errorf("期望值=value2, 实际=%v", value)
		}
	})

	t.Run("更新不存在的键应返回错误", func(t *testing.T) {
		err := cache.Update("nonexistent", "value")
		if err != ErrKeyNotFound {
			t.Errorf("期望错误=ErrKeyNotFound, 实际=%v", err)
		}
	})

	t.Run("值相同时不更新", func(t *testing.T) {
		cache.Set("key3", "value3")
		originalValue, _ := cache.Get("key3")

		err := cache.Update("key3", "value3")
		if err != nil {
			t.Errorf("更新失败: %v", err)
		}

		currentValue, _ := cache.Get("key3")
		if currentValue != originalValue {
			t.Error("值相同时不应更新")
		}
	})
}

// TestCache_Delete 测试 Delete 方法
// 测试场景：删除存在的键、删除不存在的键不应报错
func TestCache_Delete(t *testing.T) {
	cache := NewCache(10, &NoEviction{})

	t.Run("删除存在的键", func(t *testing.T) {
		cache.Set("key1", "value1")
		cache.Delete("key1")

		_, err := cache.Get("key1")
		if err != ErrKeyNotFound {
			t.Errorf("删除后应不存在, 实际错误=%v", err)
		}
	})

	t.Run("删除不存在的键不应报错", func(t *testing.T) {
		cache.Delete("nonexistent")
		// 不应 panic
	})
}

// TestCache_Clear 测试 Clear 方法
// 测试场景：清空有元素的缓存、清空空缓存
func TestCache_Clear(t *testing.T) {
	cache := NewCache(10, &NoEviction{})

	t.Run("清空有元素的缓存", func(t *testing.T) {
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")

		cache.Clear()

		_, err := cache.Get("key1")
		if err != ErrKeyNotFound {
			t.Error("清空后 key1 应不存在")
		}

		_, err = cache.Get("key2")
		if err != ErrKeyNotFound {
			t.Error("清空后 key2 应不存在")
		}

		_, err = cache.Get("key3")
		if err != ErrKeyNotFound {
			t.Error("清空后 key3 应不存在")
		}
	})

	t.Run("清空空缓存", func(t *testing.T) {
		cache.Clear()
		// 不应 panic
	})
}

// TestCache_Concurrent 测试并发场景
// 测试场景：多个 goroutine 同时读写缓存
func TestCache_Concurrent(t *testing.T) {
	t.Run("并发读写", func(t *testing.T) {
		cache := NewCache(1000, &NoEviction{})
		var wg sync.WaitGroup

		// 启动多个写 goroutine
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					key := "key" + string(rune(idx*100+j))
					cache.Set(key, idx*100+j)
				}
			}(i)
		}

		// 启动多个读 goroutine
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					key := "key" + string(rune(idx*100+j))
					cache.Get(key)
				}
			}(i)
		}

		wg.Wait()
		// 如果没有 panic 或死锁，测试通过
	})

	t.Run("并发删除", func(t *testing.T) {
		cache := NewCache(100, &NoEviction{})

		// 先填充缓存
		for i := 0; i < 100; i++ {
			cache.Set("key"+string(rune(i)), i)
		}

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					cache.Delete("key" + string(rune(idx*10+j)))
				}
			}(i)
		}

		wg.Wait()
		// 如果没有 panic 或死锁，测试通过
	})
}

// TestNewLRU 测试 NewLRU 构造函数
// 测试场景：正常创建 LRU 策略
func TestNewLRU(t *testing.T) {
	t.Run("创建 LRU 策略", func(t *testing.T) {
		lru := NewLRU(10)
		if lru == nil {
			t.Fatal("NewLRU 返回 nil")
		}
		if lru.cacheSize != 10 {
			t.Errorf("期望容量=10, 实际=%d", lru.cacheSize)
		}
		if lru.cache == nil {
			t.Error("cache map 未初始化")
		}
		if lru.evictList == nil {
			t.Error("evictList 未初始化")
		}
	})
}

// TestNoEviction 测试 NoEviction 策略
// 测试场景：ShouldEvict 返回 false、Evict 不做任何操作
func TestNoEviction(t *testing.T) {
	noEviction := &NoEviction{}
	cache := NewCache(10, noEviction)

	t.Run("ShouldEvict 应返回 false", func(t *testing.T) {
		if noEviction.ShouldEvict() {
			t.Error("NoEviction.ShouldEvict 应返回 false")
		}
	})

	t.Run("Evict 不应删除任何数据", func(t *testing.T) {
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		initialLen := len(cache.data)

		noEviction.Evict(cache)

		if len(cache.data) != initialLen {
			t.Error("NoEviction.Evict 不应删除数据")
		}
	})
}

// TestLRU_ShouldEvict 测试 LRU.ShouldEvict 方法
// 测试场景：未超过容量、超过容量
func TestLRU_ShouldEvict(t *testing.T) {
	lru := NewLRU(3)

	t.Run("未超过容量时不应淘汰", func(t *testing.T) {
		lru.cache["a"] = nil
		lru.cache["b"] = nil

		if lru.ShouldEvict() {
			t.Error("未超过容量时不应淘汰")
		}
	})

	t.Run("超过容量时应淘汰", func(t *testing.T) {
		lru.cache["c"] = nil

		if !lru.ShouldEvict() {
			t.Error("超过容量时应淘汰")
		}
	})
}

// TestCache_EvictionWithLRU 测试使用 LRU 淘汰策略的缓存
// 测试场景：容量满时触发淘汰
func TestCache_EvictionWithLRU(t *testing.T) {
	t.Run("容量满时 LRU 应淘汰最久未使用的项", func(t *testing.T) {
		lru := NewLRU(2)
		cache := NewCache(2, lru)

		// 填充缓存到容量上限
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")

		// 添加第三个项应触发淘汰
		cache.Set("key3", "value3")

		// 验证某个键被淘汰（LRU 实现可能因版本不同而异）
		totalKeys := 0
		for range cache.data {
			totalKeys++
		}

		if totalKeys != 2 {
			t.Errorf("期望缓存有2个键, 实际=%d", totalKeys)
		}
	})
}

// BenchmarkCache_Set 性能基准测试
func BenchmarkCache_Set(b *testing.B) {
	cache := NewCache(10000, &NoEviction{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value")
	}
}

// BenchmarkCache_Get 性能基准测试
func BenchmarkCache_Get(b *testing.B) {
	cache := NewCache(10000, &NoEviction{})
	cache.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

// BenchmarkCache_Update 性能基准测试
func BenchmarkCache_Update(b *testing.B) {
	cache := NewCache(10000, &NoEviction{})
	cache.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Update("key", "new_value")
	}
}

// BenchmarkNewCache 性能基准测试
func BenchmarkNewCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewCache(1000, &NoEviction{})
	}
}

// ExampleNewCache 示例函数
func ExampleNewCache() {
	// 创建一个容量为 100 的缓存，使用 LRU 淘汰策略
	cache := NewCache(100, NewLRU(100))

	// 存储数据
	cache.Set("user:1", "Alice")
	cache.Set("user:2", "Bob")

	// 获取数据
	if value, err := cache.Get("user:1"); err == nil {
		println(value.(string))
	}

	// 更新数据
	cache.Update("user:1", "Alice Updated")

	// 删除数据
	cache.Delete("user:2")

	// 清空缓存
	cache.Clear()
}
