package shortener

import (
	"URL_Shortener/config"
	"URL_Shortener/pkg/utils/common"
	"URL_Shortener/pkg/utils/trace"
	"context"
	"fmt"
	"time"
)

type murmurShortener struct {
	hashPool *common.HashPool
	keyPool  chan string
}

type MurMurShortenerConfig struct {
	HashPoolSize int
}

func newMurmurShortener(conf MurMurShortenerConfig) Shortener {
	murmurShortener := &murmurShortener{
		hashPool: common.NewHashPool(conf.HashPoolSize),
		keyPool:  make(chan string, conf.HashPoolSize),
	}
	go murmurShortener.GenerateUrlId()
	return murmurShortener
}
func (s *murmurShortener) GenerateUrlId() (string, error) {

	for {
		hash := s.hashPool.GetHash()
		key := fmt.Sprintf("%s", time.Now().Format(time.RFC3339Nano))
		_, err := hash.Write([]byte(key))
		if err != nil {
			return "", err
		}

		murmurHash := hash.Sum32()
		s.keyPool <- common.Base62Encode(murmurHash)
		s.hashPool.ReleaseHash(hash)
	}
}
func (s *murmurShortener) GetUrlId(ctx context.Context, url string) (string, error) {
	if config.GetConfig().Trace.Enable {
		c, span := trace.NewSpan(ctx, "http://jaeger:14268/api/traces")
		defer span.End()
		ctx = c
	}

	return <-s.keyPool, nil
}

func (s *murmurShortener) Close() {
	s.hashPool.Close()
}
