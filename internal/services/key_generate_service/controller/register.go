package controller

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/services/key_generate_service/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	handler, err := handler.NewKeyHandler(handler.DefaultKeyHandlerConf{
		HashPoolSize:   config.GetConfig().HashPoolSize,
		StoreBatchSize: config.GetConfig().StoreBatchSize,
	})
	if err != nil {
		panic(err)
	}
	s := NewController(handler)

	r.POST("/api/v1/key", s.GenerateKey)

}
