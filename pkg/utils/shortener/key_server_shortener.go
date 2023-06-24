package shortener

import (
	"URL_Shortener/internal/models"
	"URL_Shortener/pkg/utils/common"
	"URL_Shortener/pkg/utils/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type KeyServerShortenerConfig struct {
	KeyServerAddr string
	KeyPoolSize   int
	RetryTimes    int
}

type keyServerShortener struct {
	wait          bool
	cond          *sync.Cond
	keyServerAddr string
	keyPool       chan string
	KeyPoolSize   int
	RetryTimes    int
}

func newKeyServerShortener(conf KeyServerShortenerConfig) Shortener {
	keyServerShortener := keyServerShortener{
		keyServerAddr: conf.KeyServerAddr,
		keyPool:       make(chan string, conf.KeyPoolSize),
		KeyPoolSize:   conf.KeyPoolSize,
		RetryTimes:    conf.RetryTimes,
	}
	go keyServerShortener.GetNewKey()
	return &keyServerShortener
}

func (s *keyServerShortener) GenerateUrlId(url string) (string, error) {

	for {
		select {
		case code, ok := <-s.keyPool:
			if !ok {
				s.cond.L.Unlock()
				return "", fmt.Errorf("GenerateCode: key pool is closed")
			}
			return code, nil
		case <-time.After(10 * time.Second):
			return "", fmt.Errorf("GenerateCode: timeout")
		}
	}
}

func (s *keyServerShortener) GetNewKey() {

	fibonacci := common.NewFibonacci(1 * time.Second)

	retryCount := 0
	for {
		if retryCount == s.RetryTimes {
			fibonacci = common.NewFibonacci(1 * time.Second)
		}

		var response *http.Response

		url := fmt.Sprintf("%s/api/v1/key", s.keyServerAddr)
		bs, _ := json.Marshal(models.GetKeysRequest{
			Nums: s.KeyPoolSize,
		})
		reader := bytes.NewReader(bs)
		var err error

		req, err := http.NewRequest(http.MethodPost, url, reader)
		req.Header.Set("Content-Type", "application/json")

		if response, err = http.DefaultClient.Do(req); err != nil {
			fmt.Println(err)
			retryWait, _ := fibonacci.Next()
			t := time.NewTimer(retryWait)
			<-t.C
			retryCount++
			continue
		}

		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			logger.LoadExtra(map[string]interface{}{
				"status_code": response.StatusCode,
			}).Error("GetNewKey: response status code is not 200")
			retryCount++
			continue
		}
		defer response.Body.Close()
		bs, err = io.ReadAll(response.Body)
		result := models.GetKeysResponse{}
		if err = json.Unmarshal(bs, &result); err != nil {
			logger.LoadExtra(map[string]interface{}{
				"error": err,
			}).Error("GetNewKey: Unmarshal response error")
			retryCount++
			continue
		}
		for _, key := range result.Keys {
			s.keyPool <- key
		}

	}
}

func (s *keyServerShortener) Close() {
	close(s.keyPool)
}
