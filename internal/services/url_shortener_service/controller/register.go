package controller

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/services/url_shortener_service/handler"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	handlerConf := handler.DefaultHandlerConf{
		HashPoolSize: config.GetConfig().HashPoolSize,
	}
	if config.GetConfig().EnableKeyService {
		keyServerAdd := fmt.Sprintf("%s:%s",
			config.GetConfig().EndPoints.KeyServer.Http.Host,
			config.GetConfig().EndPoints.KeyServer.Http.Port)

		if config.GetConfig().EndPoints.KeyServer.Http.EnableTls {
			keyServerAdd = "https://" + keyServerAdd
		} else {
			keyServerAdd = "http://" + keyServerAdd
		}
		handlerConf.KeyServiceAddr = keyServerAdd
		handlerConf.EnableKeyService = true
	}
	handler, err := handler.NewDefaultShortenerHandler(handlerConf)

	if err != nil {
		panic(err)
	}

	s := NewController(handler)

	r.POST("/api/v1/urls", s.NewShortUrl)
	r.GET("/:urlId", s.RedirectUrl)
}
