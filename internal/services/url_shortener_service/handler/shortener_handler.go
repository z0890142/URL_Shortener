package handler

import (
	"URL_Shortener/c"
	"URL_Shortener/config"
	"URL_Shortener/internal/data/url_mapping_data"
	"URL_Shortener/pkg/utils/common"
	"URL_Shortener/pkg/utils/logger"
	"URL_Shortener/pkg/utils/shortener"
	"URL_Shortener/pkg/utils/trace"
	"context"
	"fmt"
	"sync"
)

type defaultShortenerHandler struct {
	mu                  *sync.Mutex
	murmurShortener     shortener.Shortener
	keyShortener        shortener.Shortener
	urlMappingDataRedis url_mapping_data.UrlMappingData
	urlMappingDataMysql url_mapping_data.UrlMappingData
}

type DefaultHandlerConf struct {
	EnableKeyService bool
	KeyServiceAddr   string
	HashPoolSize     int
	RedisOpts        config.RedisOption
	DatabaseOpts     config.DatabaseOption
}

func NewDefaultShortenerHandler(conf DefaultHandlerConf) (ShortenerHandler, error) {
	murmurShortener := shortener.NewShortener(shortener.MurMurShortenerConfig{
		HashPoolSize: conf.HashPoolSize,
	})

	urlMappingDataMysql, err := url_mapping_data.NewUrlMappingData(conf.DatabaseOpts)
	if err != nil {
		return nil, fmt.Errorf("NewDefaultShortenerHandler: %w", err)
	}

	handler := defaultShortenerHandler{
		murmurShortener:     murmurShortener,
		urlMappingDataMysql: urlMappingDataMysql,
		mu:                  &sync.Mutex{},
	}

	if conf.RedisOpts.Enable {
		urlMappingDataRedis, err := url_mapping_data.NewUrlMappingData(url_mapping_data.UrlMappingRedisConfig{
			Host:     conf.RedisOpts.Host,
			Port:     conf.RedisOpts.Port,
			Password: conf.RedisOpts.Password,
		})
		if err != nil {
			logger.LoadExtra(map[string]interface{}{
				"redis": conf.RedisOpts,
				"error": err,
			}).Error("NewDefaultShortenerHandler: NewUrlMappingData")
		} else {
			handler.urlMappingDataRedis = urlMappingDataRedis
		}
	}

	if conf.EnableKeyService {
		keyShortener := shortener.NewShortener(shortener.KeyServerShortenerConfig{
			KeyServerAddr: conf.KeyServiceAddr,
			KeyPoolSize:   conf.HashPoolSize,
		})
		handler.keyShortener = keyShortener
	}
	return &handler, nil
}

func (h *defaultShortenerHandler) GenerateShortUrl(ctx context.Context, url, expireAt string) (urlId string, err error) {
	if config.GetConfig().Trace.Enable {
		c, span := trace.NewSpan(ctx, "http://jaeger:14268/api/traces")
		defer span.End()
		ctx = c
	}

	for {
		urlId, err = h.GetUrlId(ctx, url)
		if err != nil {
			return "", fmt.Errorf("GenerateShortUrl: %w", err)
		}

		//set urlId in mysql
		err = h.urlMappingDataMysql.SetUrlId(ctx, urlId, url, expireAt)
		if err != nil && common.SqlErrCode(err) == c.MySQLErrDuplicateEntryCode {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("GenerateShortUrl: %w", err)
		}

		if h.urlMappingDataRedis == nil {
			return urlId, nil
		}

		//set urlId in redis
		if err := h.urlMappingDataRedis.SetUrlId(ctx, urlId, url, expireAt); err != nil {
			logger.LoadExtra(map[string]interface{}{
				"urlId":    urlId,
				"url":      url,
				"expireAt": expireAt,
				"error":    err,
			}).Error("GenerateShortUrl: SetUrlId in redis failed")
		}
		return urlId, nil
	}

}

func (h *defaultShortenerHandler) GetUrl(ctx context.Context, urlId string) (url string, err error) {
	if config.GetConfig().Trace.Enable {
		c, span := trace.NewSpan(ctx, "http://jaeger:14268/api/traces")
		defer span.End()
		ctx = c
	}

	if h.urlMappingDataRedis == nil {
		return h.urlMappingDataMysql.GetUrl(ctx, urlId)
	}

	logger.Info("GetUrl: GetUrl from redis")
	url, err = h.urlMappingDataRedis.GetUrl(ctx, urlId)
	if err == nil {
		return url, nil
	}

	logger.Info("GetUrl: GetUrl from DB")
	url, err = h.urlMappingDataMysql.GetUrl(ctx, urlId)
	if err != nil {
		return "", fmt.Errorf("GetUrl: %w", err)
	}

	return url, nil

}

func (h *defaultShortenerHandler) GetUrlId(ctx context.Context, url string) (string, error) {
	if config.GetConfig().Trace.Enable {
		c, span := trace.NewSpan(ctx, "http://jaeger:14268/api/traces")
		defer span.End()
		ctx = c
	}

	var urlId string
	var err error
	if h.keyShortener == nil {
		urlId, err = h.murmurShortener.GetUrlId(ctx, url)
		if err != nil {
			return "", fmt.Errorf("GenerateShortUrl: %w", err)
		}
		return urlId, nil
	}

	if urlId, err = h.keyShortener.GetUrlId(ctx, url); err != nil {
		urlId, err = h.murmurShortener.GetUrlId(ctx, url)
		logger.Info("GetUrlId: GenerateUrlId from murmur")
		if err != nil {
			return "", fmt.Errorf("GenerateShortUrl: %w", err)
		}
	}
	return urlId, nil
}

func (h *defaultShortenerHandler) Shutdown() {
	h.murmurShortener.Close()
	h.keyShortener.Close()
	h.urlMappingDataRedis.Close()
	h.urlMappingDataMysql.Close()
}
