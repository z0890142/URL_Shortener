package controller

import (
	"URL_Shortener/internal/services/key_generate_service/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	handler, err := handler.NewDefaultKeyHandler()
	if err != nil {
		panic(err)
	}
	s := NewController(handler)

	r.POST("/api/v1/key", s.GenerateKey)

}
