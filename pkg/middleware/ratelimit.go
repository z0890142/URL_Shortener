package middleware

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/data/ratelimit_data"
	"URL_Shortener/pkg/utils/algorithm"
	"URL_Shortener/pkg/utils/hotp"
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
)

type RatelimitConfig struct {
	Second int
	Number int64
}

func Ratelimit(conf RatelimitConfig) gin.HandlerFunc {
	ratelimitData, err := ratelimit_data.NewRateLimitData(&ratelimit_data.RatelimitDataConfig{
		Host:     config.GetConfig().RatelimitRedis.Host,
		Port:     config.GetConfig().RatelimitRedis.Port,
		Password: config.GetConfig().RatelimitRedis.Password,
	})
	fmt.Println(err)

	return func(ctx *gin.Context) {

		counter := int64(math.Floor(float64(time.Now().UTC().Unix()) / float64(conf.Second)))
		generateOpts := hotp.GenerateOptions{
			Algorithm: algorithm.AlgorithmSHA1,
			Digits:    6,
		}
		addr := ctx.Request.RemoteAddr
		generateCode, err := hotp.HotpGenerateCode(addr, uint64(counter), generateOpts)
		fmt.Println(generateCode)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
		n, err := ratelimitData.Set(generateCode)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
		if n > conf.Number {
			ctx.AbortWithStatusJSON(429, gin.H{"error": "too many request"})
			return
		}
	}
}
