package redis

import (
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Client = redis.Client

var instance *redis.Client

func Instance() *redis.Client {
	var once sync.Once
	once.Do(func() {
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			host = "127.0.0.1"
		}

		instance = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, 6379),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	})
	return instance
}
