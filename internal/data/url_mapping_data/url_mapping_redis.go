package url_mapping_data

import (
	"URL_Shortener/config"
	"URL_Shortener/pkg/utils/trace"
	"context"
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

func (r *urlMappingRedis) SetUrlId(ctx context.Context, urlId, url, expireAt string) error {
	if config.GetConfig().Trace.Enable {
		c, span := trace.NewSpan(ctx, "http://jaeger:14268/api/traces")
		defer span.End()
		ctx = c
	}

	t, err := time.Parse(time.RFC3339, expireAt)
	duration := t.Sub(time.Now())
	seconds := int(duration.Seconds())

	if err != nil {
		return fmt.Errorf("SetUrlId: %s", err)
	}
	return r.redisClient.WithContext(ctx).Set(urlId, url, time.Duration(seconds)*time.Second).Err()
}

func (r *urlMappingRedis) GetUrl(ctx context.Context, urlId string) (string, error) {
	if config.GetConfig().Trace.Enable {
		c, span := trace.NewSpan(ctx, "http://jaeger:14268/api/traces")
		defer span.End()
		ctx = c
	}

	return r.redisClient.WithContext(ctx).Get(urlId).Result()
}

func (r *urlMappingRedis) Close() error {
	return r.redisClient.Close()
}
