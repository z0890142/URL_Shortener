package shortener

import "context"

type Shortener interface {
	GetUrlId(ctx context.Context, url string) (string, error)

	Close()
}

func NewShortener(conf interface{}) Shortener {
	switch c := conf.(type) {
	case MurMurShortenerConfig:
		return newMurmurShortener(c)
	case KeyServerShortenerConfig:
		return newKeyServerShortener(c)
	}
	return nil
}
