package url_mapping_data

import (
	"URL_Shortener/c"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type urlMappingRedis struct {
	redisClient *redis.Client
}
type UrlMappingRedisConfig struct {
	Host     string
	Port     string
	Password string
}

func newUrlMappingRedis(conf UrlMappingRedisConfig) (UrlMappingData, error) {
	addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
	})
	return &urlMappingRedis{redisClient: cli}, nil
}

func (r *urlMappingRedis) SetUrlId(urlId, url, expireAt string) error {
	expiration, err := time.Parse(c.TimeFormat, expireAt)
	duration := expiration.Sub(time.Now())
	if err != nil {
		return fmt.Errorf("SetUrlId: %s", err)
	}
	return r.redisClient.Set(urlId, url, duration).Err()
}

func (r *urlMappingRedis) GetUrl(urlId string) (string, error) {
	return r.redisClient.Get(urlId).Result()
}

func (r *urlMappingRedis) Close() error {
	return r.redisClient.Close()
}
