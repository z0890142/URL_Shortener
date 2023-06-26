package handler

import "context"

type ShortenerHandler interface {
	GenerateShortUrl(ctx context.Context, url, expireAt string) (string, error)
	GetUrl(ctx context.Context, urlId string) (string, error)
	Shutdown()
}
