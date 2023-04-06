package cache

import (
	"github.com/patrickmn/go-cache"
)

const ContextSuffix = "_content"

type MemoryStore struct {
	cache *cache.Cache // 本地缓存
}

func (m *MemoryStore) SetChatContext(userId string, content string) {
	m.cache.Set(userId+ContextSuffix, content, cache.DefaultExpiration)
}

func (m *MemoryStore) GetChatContext(userId string) string {
	if context, ok := m.cache.Get(userId + ContextSuffix); ok {
		return context.(string)
	}
	return ""
}
