package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	config     *Config
	client     *redis.Client
	connected  bool
	maxRetries int
}

func NewRedis(config *Config) *Redis {
	return &Redis{
		config:     config,
		maxRetries: 3,
	}
}

func (r *Redis) Connect(ctx context.Context) error {
	r.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.config.Host, r.config.Port),
		Password: r.config.Password,
		DB:       0,
	})

	err := r.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	r.connected = true
	return nil
}

func (r *Redis) Disconnect(ctx context.Context) error {
	if r.client != nil {
		err := r.client.Close()
		if err != nil {
			return fmt.Errorf("failed to disconnect from Redis: %v", err)
		}
		r.connected = false
	}
	return nil
}

func (r *Redis) Ping(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("database not connected")
	}
	return r.client.Ping(ctx).Err()
}

func (r *Redis) IsConnected() bool {
	return r.connected
}

func (r *Redis) Reconnect(ctx context.Context) error {
	r.Disconnect(ctx)

	for i := 0; i < r.maxRetries; i++ {
		err := r.Connect(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return fmt.Errorf("failed to reconnect after %d attempts", r.maxRetries)
}

// Redis specific operations
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}
