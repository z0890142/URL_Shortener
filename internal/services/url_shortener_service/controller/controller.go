package controller

import (
	"URL_Shortener/internal/models"
	"URL_Shortener/internal/services/url_shortener_service/handler"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type ShortenerController struct {
	shortHandler  handler.ShortenerHandler
	shuntDownOnce sync.Once
}

func NewController(shortHandler handler.ShortenerHandler) *ShortenerController {
	return &ShortenerController{
		shortHandler:  shortHandler,
		shuntDownOnce: sync.Once{},
	}
}

var NewShortUrlReqPool = sync.Pool{
	New: func() interface{} {
		return new(models.NewShortUrlRequest)
	},
}

func (s *ShortenerController) NewShortUrl(c *gin.Context) {
	req := NewShortUrlReqPool.Get().(*models.NewShortUrlRequest)
	defer NewShortUrlReqPool.Put(req)

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := validReq(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := s.shortHandler.GenerateShortUrl(req.Url, req.ExpireAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"shortUrl": fmt.Sprintf("http://localhost/%s", id),
	})

}

func (s *ShortenerController) RedirectUrl(c *gin.Context) {

	id := c.Param("urlId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "urlId is empty"})
		return
	}

	url, err := s.shortHandler.GetUrl(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}

func (s *ShortenerController) Shutdown() {
	s.shuntDownOnce.Do(func() {
		s.shortHandler.Shutdown()
	})
}
