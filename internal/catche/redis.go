package catche

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})
}

func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return Client.Set(ctx, key, data, ttl).Err()
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
