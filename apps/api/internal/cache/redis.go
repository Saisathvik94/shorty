package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		rdb: rdb,
	}
}

func NewRedisClient(connString string) (*redis.Client, error) {
	opts, err := redis.ParseURL(connString)

	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	return rdb, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", redis.Nil
	}

	if err != nil {
		return "", err
	}

	return value, nil

}
func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := c.rdb.Set(ctx, key, value, ttl).Err()

	if err != nil {
		return err
	}

	return nil

}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	_, err := c.rdb.Del(ctx, key).Result()

	if err != nil {
		return err
	}

	return nil
}
