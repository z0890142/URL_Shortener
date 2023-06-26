package controller

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/models"
	"URL_Shortener/internal/services/url_shortener_service/handler"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
	requestCtx := c.Request.Context()
	if config.GetConfig().Trace.Enable {
		span := trace.SpanFromContext(
			otel.GetTextMapPropagator().
				Extract(
					requestCtx,
					propagation.HeaderCarrier(c.Request.Header)))
		defer span.End()
	}

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

	id, err := s.shortHandler.GenerateShortUrl(requestCtx, req.Url, req.ExpireAt)
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

	url, err := s.shortHandler.GetUrl(c, id)
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
