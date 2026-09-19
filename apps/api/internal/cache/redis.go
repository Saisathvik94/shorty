package cache

import "github.com/redis/go-redis/v9"

func NewRedisClient(connString string) (*redis.Client, error) {
	opts, err := redis.ParseURL(connString)

	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	return rdb, nil
}
