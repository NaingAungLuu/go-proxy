package cache

type CacheStore interface {
	Put(key, value string)
	Get(key string) string
}

type OnMemoryCache struct {
	memory map[string]string
}

func (s *OnMemoryCache) Put(key, value string) {

}

func (s *OnMemoryCache) Get(key string) string {
	return "Hello, World"
}

func NewCacheStore() CacheStore {
	return &OnMemoryCache{}
}
