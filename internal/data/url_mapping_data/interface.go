package url_mapping_data

import (
	"URL_Shortener/config"
	"context"
	"fmt"
)

type UrlMappingData interface {
	SetUrlId(ctx context.Context, urlId, url, expireAt string) error
	GetUrl(ctx context.Context, urlId string) (string, error)
	Close() error
}

func NewUrlMappingData(conf interface{}) (UrlMappingData, error) {
	switch c := conf.(type) {
	case config.DatabaseOption:
		return newUrlMappingMysql(c)
	case UrlMappingRedisConfig:
		return newUrlMappingRedis(c)
	}
	return nil, fmt.Errorf("NewUrlMappingData: unknown config type")
}
