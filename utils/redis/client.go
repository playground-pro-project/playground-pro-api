package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
	log    *zap.Logger
}

func NewRedisClient() *RedisClient {
	log := middlewares.Log()
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "172.31.3.99:6379",
		Password: "",
		DB:       0,
	})

	return &RedisClient{
		client: client,
		ctx:    ctx,
		log:    log,
	}
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) SetOTP(key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(r.ctx, key, value, expiration).Err()
	if err != nil {
		r.log.Error("Failed to set OTP in Redis", zap.Error(err))
		return fmt.Errorf("failed to set OTP in Redis: %w", err)
	}

	return nil
}

func (r *RedisClient) GetOTP(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("OTP not found for key: %s", key)
	} else if err != nil {
		r.log.Error("Failed to get OTP from Redis", zap.Error(err))
		return "", fmt.Errorf("failed to get OTP from Redis: %w", err)
	}

	return val, nil
}
