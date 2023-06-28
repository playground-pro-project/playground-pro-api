package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/venue"
)

var (
	log         = middlewares.Log()
	redisClient *redis.Client
)

func InitRedis(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		// Addr:     config.REDIS_HOST + ":" + config.REDIS_PORT,
		// Password: config.REDIS_PASSWORD,
		// DB:       config.REDIS_DATABASE,

		Addr:     "172.31.3.99:6379",
		Password: "",
		DB:       0,

		// Addr:     "localhost:6379",
		// Password: "",
		// DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Sugar().Error(err)
	}
	log.Info(pong)

	return client, nil
}

// RedisClient returns the Redis client instance
func RedisClient(ctx context.Context) (*redis.Client, error) {
	// Initialize Redis client if not already done
	if redisClient == nil {
		var err error
		redisClient, err = InitRedis(ctx)
		if err != nil {
			return nil, err
		}
	}

	return redisClient, nil
}

// GetCached retrieves data from Redis cache based on the cacheKey.
func GetCached(ctx context.Context, cacheKey string) (interface{}, error) {
	if redisClient == nil {
		var err error
		redisClient, err = InitRedis(ctx)
		if err != nil {
			return nil, err
		}
	}

	cachedResult, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		log.Sugar().Error("error while retrieving data from Redis cache:", err)
		return nil, err
	} else if cachedResult != "" {
		var result []venue.VenueCore
		err = json.Unmarshal([]byte(cachedResult), &result)
		if err != nil {
			log.Sugar().Error("error while unmarshaling cached result:", err)
			return nil, err
		} else {
			log.Sugar().Info("venue data found in Redis cache")
			return result, nil
		}
	}

	return nil, nil
}

// SetCached sets the provided venues data into Redis cache.
func SetCached(ctx context.Context, cacheKey string, value interface{}, expTime time.Duration) error {
	if redisClient == nil {
		var err error
		redisClient, err = InitRedis(ctx)
		if err != nil {
			return err
		}
	}

	resultBytes, err := json.Marshal(value)
	if err != nil {
		log.Sugar().Error("error marshaling result:", err)
		return err
	}

	err = redisClient.Set(ctx, cacheKey, string(resultBytes), expTime).Err()
	if err != nil {
		log.Sugar().Error("error while setting cache:", err)
	}

	return err
}
