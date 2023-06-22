package handler

import (
	"URL_Shortener/c"
	"URL_Shortener/config"
	"URL_Shortener/internal/data/url_mapping_data"
	"URL_Shortener/internal/utils/logger"
	"URL_Shortener/internal/utils/shortener"
	"fmt"
)

type defaultShortenerHandler struct {
	murmurShortener     shortener.Shortener
	keyShortener        shortener.Shortener
	urlMappingDataRedis url_mapping_data.UrlMappingData
	urlMappingDataMysql url_mapping_data.UrlMappingData
}

type DefaultHandlerConf struct {
	EnableKeyService bool
	KeyServiceAddr   string
	HashPoolSize     int
}

func NewDefaultShortenerHandler(conf DefaultHandlerConf) (ShortenerHandler, error) {
	murmurShortener := shortener.NewShortener(shortener.MurMurShortenerConfig{
		HashPoolSize: conf.HashPoolSize,
	})

	urlMappingDataMysql, err := url_mapping_data.NewUrlMappingData(config.GetConfig().Databases)
	if err != nil {
		return nil, fmt.Errorf("NewDefaultShortenerHandler: %w", err)
	}

	urlMappingDataRedis, err := url_mapping_data.NewUrlMappingData(url_mapping_data.UrlMappingRedisConfig{
		Host:     config.GetConfig().Redis.Host,
		Port:     config.GetConfig().Redis.Port,
		Password: config.GetConfig().Redis.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("NewDefaultShortenerHandler: %w", err)
	}

	handler := defaultShortenerHandler{
		murmurShortener:     murmurShortener,
		urlMappingDataRedis: urlMappingDataRedis,
		urlMappingDataMysql: urlMappingDataMysql,
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

func (h *defaultShortenerHandler) GenerateShortUrl(url, expireAt string) (urlId string, err error) {

	for {
		urlId, err = h.GetUrlId(url)
		if err != nil {
			return "", fmt.Errorf("GenerateShortUrl: %w", err)
		}

		//set urlId in mysql
		err = h.urlMappingDataMysql.SetUrlId(urlId, url, expireAt)
		if err != nil && SqlErrCode(err) == c.MySQLErrDuplicateEntryCode {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("GenerateShortUrl: %w", err)
		}

		//set urlId in redis
		if err := h.urlMappingDataRedis.SetUrlId(urlId, url, expireAt); err != nil {
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

func (h *defaultShortenerHandler) GetUrl(urlId string) (url string, err error) {
	if url, err = h.urlMappingDataRedis.GetUrl(urlId); err != nil {
		url, err = h.urlMappingDataMysql.GetUrl(urlId)
		if err != nil {
			return "", fmt.Errorf("GetUrl: %w", err)
		}
	}
	return url, nil

}

func (h *defaultShortenerHandler) GetUrlId(url string) (string, error) {
	var urlId string
	var err error
	if h.keyShortener == nil {
		urlId, err = h.murmurShortener.GenerateUrlId(url)
		if err != nil {
			return "", fmt.Errorf("GenerateShortUrl: %w", err)
		}
		return urlId, nil
	}
	if urlId, err = h.keyShortener.GenerateUrlId(url); err != nil {
		urlId, err = h.murmurShortener.GenerateUrlId(url)
		if err != nil {
			return "", fmt.Errorf("GenerateShortUrl: %w", err)
		}
	}
	return urlId, nil
}
