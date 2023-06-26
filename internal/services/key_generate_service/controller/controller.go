package controller

import (
	"URL_Shortener/internal/models"
	"URL_Shortener/internal/services/key_generate_service/handler"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type KeyController struct {
	keyHandler    handler.KeyHandler
	shuntDownOnce sync.Once
}

func NewController(keyHandler handler.KeyHandler) *KeyController {
	return &KeyController{
		keyHandler:    keyHandler,
		shuntDownOnce: sync.Once{},
	}

}

var NewGetKeysReqPool = sync.Pool{
	New: func() interface{} {
		return new(models.GetKeysRequest)
	},
}

func (s *KeyController) GenerateKey(c *gin.Context) {
	req := NewGetKeysReqPool.Get().(*models.GetKeysRequest)
	defer NewGetKeysReqPool.Put(req)

	if err := c.ShouldBindJSON(req); err != nil {
		handleError(c, http.StatusBadRequest, err)
		return
	}

	err := validReq(req)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err)
		return
	}

	keys, err := s.keyHandler.GetKeys(req.Nums)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, models.GetKeysResponse{
		Keys: keys,
	})

}

func (s *KeyController) Shutdown() {
	s.shuntDownOnce.Do(func() {
		s.keyHandler.Shutdown()
	})
}

func handleError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}
