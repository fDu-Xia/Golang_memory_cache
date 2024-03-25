package cache

import (
	utils2 "memory-cache/cache/utils"
	"sync"
	"time"
)

type Value struct {
	val        interface{}
	expireTime time.Time
}

type MemoryCache struct {
	//单位为bytes
	curSize    int64
	maxSize    int64
	maxSizeStr string
	rwMutex    sync.RWMutex
	values     map[string]*Value
}

func (m *MemoryCache) get(key string) (interface{}, bool) {
	v, exist := m.values[key]
	return v, exist
}

func (m *MemoryCache) del(key string) {
	delete(m.values, key)
}

func NewMemoryCache(maxSize string) *MemoryCache {
	m := &MemoryCache{
		values: make(map[string]*Value),
	}
	m.SetMaxMemory(maxSize)
	return m
}

func (m *MemoryCache) SetMaxMemory(size string) bool {
	m.maxSize = utils2.ParseSize(size)
	m.maxSizeStr = size
	return true
}

func (m *MemoryCache) Set(key string, val interface{}, expire ...time.Duration) bool {
	var exp time.Duration
	if len(expire) == 0 {
		exp = 0
	}
	exp = expire[0]
	v := Value{val, time.Now().Add(exp)}
	m.rwMutex.Lock()
	if _, exist := m.get(key); exist {
		m.del(key)
	}
	m.values[key] = &v
	if m.curSize+GetValSize(v) > m.maxSize {
		//这里可以完善一下，通过一些内存淘汰策略来选择删除一些 key，来判断是否还会超过最大内存
		m.del(key)
		panic("Maximum memory exceeded")
	}
	m.rwMutex.Unlock()
	return true
}

func (m *MemoryCache) Get(key string) (interface{}, bool) {
	m.rwMutex.RLock()
	v, ok := m.values[key]
	if ok {
		if v.expireTime.Before(time.Now()) {
			m.del(key)
			return nil, false
		}
	}
	m.rwMutex.RUnlock()
	return v.val, ok
}

func (m *MemoryCache) Del(key string) bool {
	m.rwMutex.Lock()
	delete(m.values, key)
	m.rwMutex.Unlock()
	return true
}

func (m *MemoryCache) Exists(key string) bool {
	m.rwMutex.RLock()
	_, ok := m.values[key]
	m.rwMutex.RUnlock()
	return ok
}

func (m *MemoryCache) Flush() bool {
	m.rwMutex.Lock()
	m.values = make(map[string]*Value)
	m.rwMutex.Unlock()
	return true
}

func (m *MemoryCache) Keys() int64 {
	defer m.rwMutex.RUnlock()
	m.rwMutex.RLock()
	return int64(len(m.values))
}

func(m *MemoryCache) cleanExpiredItem() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for k, v := range m.values {
				if v.expireTime.Before(time.Now()) {
					m.rwMutex.Lock()
					m.Del(k)
					m.rwMutex.Unlock()
				}
			}
		}
	}
}


