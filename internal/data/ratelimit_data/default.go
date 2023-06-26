package ratelimit_data

import (
	"fmt"

	"github.com/go-redis/redis"
)

type defaultRatelimitData struct {
	redisClient *redis.Client
}
type RatelimitDataConfig struct {
	Host     string
	Port     string
	Password string
}

func newDefaultRatelimitData(conf *RatelimitDataConfig) (RateLimitData, error) {
	addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
	})
	return &defaultRatelimitData{redisClient: cli}, nil
}

func (r *defaultRatelimitData) Set(key string) (int64, error) {
	return r.redisClient.Incr(key).Result()
}
