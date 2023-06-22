package handler

type ShortenerHandler interface {
	GenerateShortUrl(url, expireAt string) (string, error)
	GetUrl(urlId string) (string, error)
}
