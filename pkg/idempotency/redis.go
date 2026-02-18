package idempotency

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Store interface {
	Exists(ctx context.Context, key string) (bool, error)
	Store(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) Store {
	return &redisStore{client: client}
}

func (r *redisStore) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *redisStore) Store(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.SetNX(ctx, key, value, ttl).Err()
}
