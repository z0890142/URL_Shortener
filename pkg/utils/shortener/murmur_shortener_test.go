package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMurmurShortener(t *testing.T) {
	mur := NewShortener(&MurMurShortenerConfig{
		HashPoolSize: 0,
	})
	assert.Nil(t, mur)

	mur = NewShortener(&MurMurShortenerConfig{
		HashPoolSize: 10,
	})
	assert.NotNil(t, mur)
	for i := 0; i < 10; i++ {
		urlId, err := mur.GetUrlId(nil)
		assert.Nil(t, err)
		assert.NotNil(t, urlId)
	}
}
