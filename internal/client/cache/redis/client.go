package redis

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/katyafirstova/auth_service/internal/client/cache"
	"github.com/katyafirstova/auth_service/internal/config"
)

var _ cache.RedisClient = (*client)(nil)

type client struct {
	pool   *redis.Pool
	config config.RedisConfig
}

func NewClient(pool *redis.Pool, config config.RedisConfig) *client {
	return &client{
		pool:   pool,
		config: config,
	}
}

func (c client) Set(ctx context.Context, key string, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c client) HashSet(ctx context.Context, key string, values interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c client) HGetAll(ctx context.Context, key string) ([]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c client) Get(ctx context.Context, key string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (c client) Ping(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
