package cache

type CacheStore interface {
	Put(key, value string)
	Get(key string) string
}

type OnMemoryCache struct {
	memory map[string]string
}

func (s *OnMemoryCache) Put(key, value string) {
	s.memory[key] = value
}

func (s *OnMemoryCache) Get(key string) string {
	return s.memory[key]
}

func NewCacheStore() CacheStore {
	return &OnMemoryCache{
		memory: make(map[string]string),
	}
}

