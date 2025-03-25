package vb

import (
	"errors"
	"fmt"
	"github.piwriw.go-tools/pkg/cache"
	"reflect"
	"sync"
)

type alertCacheClient struct {
	cache *cache.Cache
	meta  map[string]any
	mu    sync.RWMutex // 使用读写锁保护并发访问
	err   []error
}

func newAlertCacheClient() *alertCacheClient {
	return &alertCacheClient{
		cache: cache.NewCache(100, &cache.NoEviction{}), // 设置缓存容量和淘汰策略
		meta:  make(map[string]any),
	}
}

func (a *alertCacheClient) GetAll() map[string]any {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cache.GetAll()
}

// GetByIndexKey 从缓存中获取数据并赋值给 result
func (a *alertCacheClient) GetByIndexKey(indexKey string, result any) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	val, ok := a.meta[indexKey]
	if !ok {
		return fmt.Errorf("metas indexKey '%s' not found", indexKey)
	}

	reflectValue := reflect.ValueOf(result)
	if reflectValue.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer")
	}

	reflectValue = reflectValue.Elem()
	if !reflectValue.CanSet() {
		return fmt.Errorf("cannot set value to result")
	}

	valValue := reflect.ValueOf(val)
	if valValue.Type() != reflectValue.Type() {
		return fmt.Errorf("type mismatch: expected %v, got %v", reflectValue.Type(), valValue.Type())
	}

	reflectValue.Set(valValue)
	return nil
}

// Table 设置当前操作的缓存表
func (a *alertCacheClient) Table(table string) *alertCacheClient {
	a.mu.Lock()
	defer a.mu.Unlock()

	get, err := a.cache.Get(table)
	if err != nil {
		a.err = append(a.err, err)
		return a
	}

	alertCacheMeta, ok := get.(map[string]any)
	if !ok {
		a.err = append(a.err, fmt.Errorf("unexpected type for table '%s'", table))
		return a
	}

	a.meta = alertCacheMeta
	return a
}

// Update 更新缓存中的数据
func (a *alertCacheClient) Update(table string, keyIndex string, data any) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	get, err := a.cache.Get(table)
	if errors.Is(err, cache.ErrKeyNotFound) {
		// 如果缓存不存在，初始化新的 map
		a.cache.Set(table, map[string]any{keyIndex: data})
		return nil
	} else if err != nil {
		return err
	}

	// 类型断言
	alertCacheMeta, ok := get.(map[string]any)
	if !ok {
		// 如果类型不匹配，初始化新的 map
		alertCacheMeta = map[string]any{keyIndex: data}
		a.cache.Set(table, alertCacheMeta)
		return nil
	}

	// 更新缓存数据
	alertCacheMeta[keyIndex] = data
	a.cache.Set(table, alertCacheMeta)
	return nil
}

// Errors 返回所有错误
func (a *alertCacheClient) Errors() []error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.err
}
