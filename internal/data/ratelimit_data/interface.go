package ratelimit_data

import "fmt"

type RateLimitData interface {
	Set(key string) (int64, error)
}

func NewRateLimitData(conf interface{}) (RateLimitData, error) {
	switch c := conf.(type) {
	case *RatelimitDataConfig:
		return newDefaultRatelimitData(c)
	}
	return nil, fmt.Errorf("unknow config type")
}
