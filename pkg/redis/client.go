package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewClientWithRetry(addr string, maxRetries int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	var err error
	for i := 0; i < maxRetries; i++ {
		err = client.Ping(context.Background()).Err()
		if err == nil {
			return client, nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return nil, err
}
