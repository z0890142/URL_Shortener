package shortener

type Shortener interface {
	GenerateUrlId(url string) (string, error)
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
