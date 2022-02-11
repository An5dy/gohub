package cache

import (
	"gohub/pkg/config"
	"gohub/pkg/redis"
	"time"
)

// RedisStore 实现 cache.Store interface
type RedisStore struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

var _ Store = (*RedisStore)(nil)

func NewRedisStore(address, username, password string, db int) *RedisStore {
	rs := &RedisStore{}
	rs.RedisClient = redis.NewClient(address, username, password, db)
	rs.KeyPrefix = config.GetString("app.name") + ":cache"
	return rs
}

func (s *RedisStore) Set(key, value string, expireTime time.Duration) {
	s.RedisClient.Set(s.cacheKey(key), value, expireTime)
}

func (s *RedisStore) Get(key string) string {
	return s.RedisClient.Get(s.cacheKey(key))
}

func (s *RedisStore) Has(key string) bool {
	return s.RedisClient.Has(s.cacheKey(key))
}

func (s *RedisStore) Forget(key string) {
	s.RedisClient.Del(s.cacheKey(key))
}

func (s *RedisStore) Forever(key, value string) {
	s.RedisClient.Set(key, value, 0)
}

func (s *RedisStore) Flush() {
	s.RedisClient.FlushDB()
}

func (s *RedisStore) IsAlive() error {
	return s.RedisClient.Ping()
}

func (s *RedisStore) Increment(parameters ...interface{}) {
	s.RedisClient.Increment(parameters...)
}

func (s *RedisStore) Decrement(parameters ...interface{}) {
	s.RedisClient.Decrement(parameters...)
}

func (s *RedisStore) cacheKey(key string) string {
	return s.KeyPrefix + key
}
