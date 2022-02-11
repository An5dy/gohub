package cache

import "time"

type Store interface {
	// Set 设置值
	Set(key, value string, expireTime time.Duration)

	// Get 获取值
	Get(key string) string

	// Has 是否设置
	Has(key string) bool

	// Forget 移除
	Forget(key string)

	// Forever 永久保存
	Forever(key, value string)

	// Flush 清空
	Flush()

	// IsAlive 是否还生效
	IsAlive() error

	// Increment 当参数只有 1 个时，为 key，增加 1。
	// 当参数有 2 个时，第一个参数为 key ，第二个参数为要增加的值 int64 类型。
	Increment(parameters ...interface{})

	// Decrement 当参数只有 1 个时，为 key，减去 1。
	// 当参数有 2 个时，第一个参数为 key ，第二个参数为要减去的值 int64 类型。
	Decrement(parameters ...interface{})
}
