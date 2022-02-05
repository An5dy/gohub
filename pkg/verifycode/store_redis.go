package verifycode

import (
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/redis"
	"time"
)

// RedisStore 实现 verifycode.Store interface
type RedisStore struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

var _ Store = (*RedisStore)(nil)

// Set 实现 verifycode.Store interface 的 Set 方法
func (s *RedisStore) Set(id, value string) bool {
	expireTime := time.Minute * time.Duration(config.GetInt64("verifycode.expire_time"))
	// 本地环境方便调试
	if app.IsLocal() {
		expireTime = time.Minute * time.Duration(config.GetInt64("verifycode.debug_expire_time"))
	}
	return s.RedisClient.Set(s.storeKey(id), value, expireTime)
}

// Get 实现 verifycode.Store interface 的 Get 方法
func (s *RedisStore) Get(id string, clear bool) string {
	key := s.storeKey(id)
	val := s.RedisClient.Get(id)
	if clear {
		s.RedisClient.Del(key)
	}
	return val
}

// Verify 实现 verifycode.Store interface 的 Verify 方法
func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v == answer
}

// storeKey 获取 redis 保存的 key
func (s *RedisStore) storeKey(key string) string {
	return s.KeyPrefix + key
}
