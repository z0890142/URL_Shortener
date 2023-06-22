package models

type NewShortUrlRequest struct {
	Url      string `json:"url"`
	ExpireAt string `son:"expireAt"`
}

type NewShortUrlResponse struct {
	Id       string `json:"id"`
	ShortUrl string `json:"shortUrl"`
}

type GetKeysRequest struct {
	Nums int `json:"nums"`
}
type GetKeysResponse struct {
	Keys []string `json:"keys"`
}
