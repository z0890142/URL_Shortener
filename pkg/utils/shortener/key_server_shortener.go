package shortener

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/models"
	"URL_Shortener/pkg/utils/logger"
	"URL_Shortener/pkg/utils/trace"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type KeyServerShortenerConfig struct {
	KeyServerAddr string
	KeyPoolSize   int
}

type keyServerShortener struct {
	keyServerAddr string
	keyPool       chan string
	KeyPoolSize   int
}

func newKeyServerShortener(conf KeyServerShortenerConfig) Shortener {
	keyServerShortener := keyServerShortener{
		keyServerAddr: conf.KeyServerAddr,
		keyPool:       make(chan string, conf.KeyPoolSize),
		KeyPoolSize:   conf.KeyPoolSize,
	}
	go keyServerShortener.GetNewKey()
	return &keyServerShortener
}

func (s *keyServerShortener) GetUrlId(ctx context.Context, url string) (string, error) {
	if config.GetConfig().Trace.Enable {
		c, span := trace.NewSpan(ctx, "http://jaeger:14268/api/traces")
		defer span.End()
		ctx = c
	}

	for {
		select {
		case code, ok := <-s.keyPool:
			if !ok {
				return "", fmt.Errorf("GenerateCode: key pool is closed")
			}
			return code, nil
		case <-time.After(10 * time.Millisecond):
			return "", fmt.Errorf("GenerateCode: timeout")
		}
	}
}

var httpClient = &http.Client{}
var NewGetKeysResponsePool = sync.Pool{
	New: func() interface{} {
		return new(models.GetKeysResponse)
	},
}

func (s *keyServerShortener) GetNewKey() {

	for {

		var response *http.Response

		url := fmt.Sprintf("%s/api/v1/key", s.keyServerAddr)
		bs, _ := json.Marshal(models.GetKeysRequest{
			Nums: s.KeyPoolSize * 2,
		})
		reader := bytes.NewReader(bs)
		var err error

		req, err := http.NewRequest(http.MethodPost, url, reader)
		req.Header.Set("Content-Type", "application/json")

		if response, err = httpClient.Do(req); err != nil {
			logger.LoadExtra(map[string]interface{}{
				"error": err,
			}).Error("GetNewKey: http request error")
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			logger.LoadExtra(map[string]interface{}{
				"status_code": response.StatusCode,
			}).Error("GetNewKey: response status code is not 200")
			continue
		}

		decoder := json.NewDecoder(response.Body)
		result := NewGetKeysResponsePool.Get().(*models.GetKeysResponse)
		if err = decoder.Decode(&result); err != nil {
			logger.LoadExtra(map[string]interface{}{
				"error": err,
			}).Error("GetNewKey: Decode response error")
			continue
		}

		for _, key := range result.Keys {
			s.keyPool <- key
		}
		NewGetKeysResponsePool.Put(result)
	}
}

func (s *keyServerShortener) Close() {
	close(s.keyPool)
}
