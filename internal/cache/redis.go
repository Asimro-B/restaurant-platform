package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	// test connection
	ctx := context.Background()
	result, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v", err)
	} else {
		log.Printf("Redis connected successfully: %s", result)
	}
}

func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	err = Client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		log.Printf("cache set error for key %s: %v", key, err)
	} else {
		log.Printf("cache set success for key: %s", key)
	}
	return err
}

func Get(ctx context.Context, key string, dest interface{}) error {
	data, err := Client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

func Delete(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

func DeleteByPattern(ctx context.Context, pattern string) error {
	keys, err := Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	return Client.Del(ctx, keys...).Err()
}
