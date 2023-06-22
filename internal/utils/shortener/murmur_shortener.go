package shortener

import (
	"URL_Shortener/internal/utils/common"
	"fmt"
	"time"
)

type murmurShortener struct {
	hashPool *common.HashPool
}

type MurMurShortenerConfig struct {
	HashPoolSize int
}

func newMurmurShortener(conf MurMurShortenerConfig) Shortener {
	return &murmurShortener{
		hashPool: common.NewHashPool(conf.HashPoolSize),
	}
}

func (s *murmurShortener) GenerateUrlId(url string) (string, error) {
	hash := s.hashPool.GetHash()
	defer s.hashPool.ReleaseHash(hash)

	key := fmt.Sprintf("%s%s", url, time.Now().Format(time.RFC3339Nano))
	_, err := hash.Write([]byte(key))
	if err != nil {
		return "", err
	}

	return common.Base62Encode(hash.Sum32()), nil
}
