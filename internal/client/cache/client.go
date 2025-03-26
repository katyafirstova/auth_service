package cache

import (
	"context"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}) error
	HashSet(ctx context.Context, key string, values interface{}) error
	HGetAll(ctx context.Context, key string) ([]interface{}, error)
	Get(ctx context.Context, key string) (interface{}, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Ping(ctx context.Context) error
}
