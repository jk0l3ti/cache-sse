package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func newRedisCache(addr string) (Cache, error) {
	redis := RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
	resp := redis.client.Ping(context.Background())
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	return &redis, nil
}

func (r *RedisCache) getKey(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// redis server needs - KEA to be enabled to send this event data to the redis client(application)
//
// redis-cli config set notify-keyspace-events KEA, and '__keyspace@0__:my-key' is the pattern
func (r *RedisCache) Stream(ctx context.Context, key string, ch chan<- any) {
	pubsub := r.client.PSubscribe(ctx, fmt.Sprintf("__keyspace@0__:%s", key))
	for msg := range pubsub.Channel() {
		switch msg.Payload {
		case "set":
			resp, err := r.getKey(ctx, key)
			if err != nil {
				fmt.Println("failed to read", key, "err: ", err.Error())
				continue
			}
			ch <- resp
		case "expired", "del":
			fmt.Println(key, "got expired or deleted")
			close(ch)
		}
	}
}
