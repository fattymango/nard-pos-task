package cache

import (
	"context"
	"multitenant/pkg/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis *redis.Client
}

func NewCache(c *config.Config) (*Cache, error) {
	r, err := NewRedisClient(c)
	if err != nil {
		return nil, err
	}
	return &Cache{redis: r}, nil
}

// Set Redis `SET key value [expiration]` command.
// Use expiration for `SETEX`-like behavior.
func (e *Cache) Set(ctx context.Context, key string, value []byte, expiration int) error {
	return e.redis.Set(ctx, key, string(value), time.Duration(expiration)*time.Second).Err()
}

// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
func (e *Cache) Get(ctx context.Context, key string) (string, error) {
	result, err := e.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// Delete Redis `DEL key` command to remove key from cache
func (e *Cache) Delete(ctx context.Context, key string) error {
	return e.redis.Del(ctx, key).Err()
}

// FlushAll Redis `FLUSHALL` command to delete all keys from the Redis database
func (e *Cache) FlushAll(ctx context.Context) error {
	return e.redis.FlushAll(ctx).Err()
}

// IncrByFloat Redis command to increment the value of a key by a float.
func (e *Cache) IncrByFloat(ctx context.Context, key string, increment float64) (float64, error) {
	return e.redis.IncrByFloat(ctx, key, increment).Result()
}

// Expire Redis command to set the expiration time of a key.
func (e *Cache) Expire(ctx context.Context, key string, expiration int) error {
	return e.redis.Expire(ctx, key, time.Duration(expiration)*time.Second).Err()
}

func (e *Cache) GetFloat(ctx context.Context, key string) (float64, error) {
	return e.redis.Get(ctx, key).Float64()
}
