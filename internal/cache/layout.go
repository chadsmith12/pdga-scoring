package cache

import "sync"

type syncedMap[K comparable, V any] struct {
    m map[K]V
    mu sync.RWMutex
}

func (m *syncedMap[K, V]) Get(key K) V {
    m.mu.RLock()
    defer m.mu.RUnlock()

    return m.m[key]
}

func (m *syncedMap[K, V]) Set(key K, value V) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.m[key] = value
}

var layoutCache = &syncedMap[int, int]{}

func InsertLayout(dbId, externalId int) {
    layoutCache.Set(externalId, dbId)
}

func LayoutDbId(externalId int) int {
   return layoutCache.Get(externalId) 
}
