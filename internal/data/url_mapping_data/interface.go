package url_mapping_data

import (
	"URL_Shortener/config"
	"fmt"
)

type UrlMappingData interface {
	SetUrlId(urlId, url, expireAt string) error
	GetUrl(urlId string) (string, error)
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
