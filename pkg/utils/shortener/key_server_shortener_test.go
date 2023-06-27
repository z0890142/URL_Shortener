package shortener

import (
	"URL_Shortener/internal/models"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestKeyServerShortener(t *testing.T) {
	router := gin.Default()
	mur := NewShortener(&MurMurShortenerConfig{
		HashPoolSize: 10,
	})
	assert.NotNil(t, mur)

	router.POST("/api/v1/key", func(c *gin.Context) {
		res := models.GetKeysRequest{}
		err := c.BindJSON(&res)
		assert.Nil(t, err)
		keys := []string{}
		for i := 0; i < res.Nums; i++ {
			urlId, err := mur.GetUrlId(nil)
			assert.Nil(t, err)
			keys = append(keys, urlId)
		}

		response := models.GetKeysResponse{
			Keys: keys,
		}
		c.JSON(200, response)
	})

	ts := httptest.NewServer(router)

	ks := NewShortener(&KeyServerShortenerConfig{
		KeyServerAddr: ts.URL,
		KeyPoolSize:   1,
		MaxRetryTimes: 10,
	})
	assert.NotNil(t, ks)

	urlId, err := ks.GetUrlId(nil)
	assert.Nil(t, err)
	assert.NotEqual(t, "", urlId)

}

func TestKeyServerShortenerFail(t *testing.T) {
	router := gin.Default()
	mur := NewShortener(&MurMurShortenerConfig{
		HashPoolSize: 10,
	})
	assert.NotNil(t, mur)

	router.POST("/api/v1/key", func(c *gin.Context) {

		c.JSON(500, nil)
	})

	ts := httptest.NewServer(router)

	ks := NewShortener(&KeyServerShortenerConfig{
		KeyServerAddr: ts.URL,
		KeyPoolSize:   1,
	})
	assert.NotNil(t, ks)

	urlId, err := ks.GetUrlId(nil)
	assert.Equal(t, "GenerateCode: timeout", err.Error())
	assert.Equal(t, "", urlId)

	ks.Close()

	urlId, err = ks.GetUrlId(nil)
	assert.Equal(t, "GenerateCode: key pool is closed", err.Error())
	assert.Equal(t, "", urlId)
}
