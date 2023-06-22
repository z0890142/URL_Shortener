package shortener

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/models"
	"URL_Shortener/internal/utils/common"
	"bytes"
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
	cond          *sync.Cond
	keyServerAddr string
	keyPool       chan string
	KeyPoolSize   int
}

func newKeyServerShortener(conf KeyServerShortenerConfig) Shortener {
	return &keyServerShortener{
		cond:          sync.NewCond(&sync.Mutex{}),
		keyServerAddr: conf.KeyServerAddr,
		keyPool:       make(chan string, conf.KeyPoolSize),
		KeyPoolSize:   conf.KeyPoolSize,
	}
}

func (s *keyServerShortener) GenerateUrlId(url string) (string, error) {

	for len(s.keyPool) == 0 {
		s.cond.Wait()
	}

	for {
		select {
		case code, ok := <-s.keyPool:
			if !ok {
				return "", fmt.Errorf("GenerateCode: key pool is closed")
			}
			if len(s.keyPool) == 0 {
				err := s.GetNewKey()
				s.cond.Broadcast()
				if err != nil {
					return "", fmt.Errorf("GenerateCode: %w", err)
				}
			}
			return code, nil
		}
	}
}

func (s *keyServerShortener) GetNewKey() error {

	url := fmt.Sprintf("%s/api/v1/key", s.keyServerAddr)
	bs, _ := json.Marshal(models.GetKeysRequest{
		Nums: s.KeyPoolSize,
	})
	body := bytes.NewBuffer(bs)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("GetNewKey: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	fibonacci := common.NewFibonacci(1 * time.Second)
	retryCount := 0

	var response *http.Response

	for {
		if retryCount > config.GetConfig().MaxRetry {
			break
		}

		if response, err = http.DefaultClient.Do(req); err != nil {
			retryCount++
			retryWait, _ := fibonacci.Next()
			t := time.NewTimer(retryWait)
			<-t.C
			break
		}
	}

	if err != nil {
		return fmt.Errorf("GetNewKey: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("GetNewKey: request error")
	}

	result := models.GetKeysResponse{}
	if err = json.Unmarshal(bs, &result); err != nil {
		return fmt.Errorf("GetNewKey: %w", err)
	}

	for _, key := range result.Keys {
		s.keyPool <- key
	}
	return nil
}
