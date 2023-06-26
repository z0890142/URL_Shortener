package key_generator

import (
	"URL_Shortener/internal/data/key_data"
	"URL_Shortener/internal/models"
	"URL_Shortener/pkg/utils/logger"
	"URL_Shortener/pkg/utils/shortener"
	"context"
	"sync"

	"github.com/google/uuid"

	"time"
)

type KeyGenerator struct {
	keyGenerationInterval time.Duration
	storeBatchSize        int
	keyData               key_data.KeyData
	murmurShortener       shortener.Shortener

	storeBuffer chan string
	stopChan    chan struct{}
}

func NewKeyGenerator(
	storeBatchSize int,
	keyData key_data.KeyData,
	murmurShortener shortener.Shortener,
) *KeyGenerator {
	return &KeyGenerator{
		keyGenerationInterval: time.Duration(1) * time.Second,
		storeBatchSize:        storeBatchSize,
		keyData:               keyData,
		murmurShortener:       murmurShortener,
		storeBuffer:           make(chan string),
		stopChan:              make(chan struct{}),
	}
}

func (k *KeyGenerator) Start() {
	go k.startRecycler(time.Second * 1)
	go k.storeKey()
	go k.generateKey()
}

func (k *KeyGenerator) startRecycler(period time.Duration) {
	ticker := time.NewTicker(period)
	windows := []int64{0, 0, 0}

	for {
		select {
		case <-ticker.C:
			total, err := k.keyData.GetAvailableKeyCount()
			if err != nil {
				logger.LoadExtra(map[string]interface{}{
					"error": err,
				}).Error("startRecycler: error")
				continue
			}

			copy(windows, windows[1:])
			windows[len(windows)-1] = total
			k.adjustKeyGenerationInterval(int(max(windows)))

		case <-k.stopChan:
			ticker.Stop()
			return
		}
	}

}

func (k *KeyGenerator) adjustKeyGenerationInterval(size int) {
	if size > k.storeBatchSize*1000 {
		if k.keyGenerationInterval == 0 {
			k.keyGenerationInterval = time.Duration(1) * time.Microsecond
		}
		k.keyGenerationInterval = k.keyGenerationInterval * 2

		if k.keyGenerationInterval > time.Duration(500)*time.Millisecond {
			k.keyGenerationInterval = time.Duration(500) * time.Millisecond
		}

	} else if size < k.storeBatchSize/2 {

		k.keyGenerationInterval = 0

	} else if size < k.storeBatchSize {

		k.keyGenerationInterval = k.keyGenerationInterval / 2
		if k.keyGenerationInterval < time.Duration(1)*time.Millisecond {
			k.keyGenerationInterval = time.Duration(1) * time.Millisecond
		}
	}
}

func (k *KeyGenerator) generateKey() {
	var uuidPool = sync.Pool{
		New: func() interface{} {
			return uuid.New()
		},
	}
	for {

		time.Sleep(k.keyGenerationInterval)
		var key string
		var err error
		uuidObj := uuidPool.Get().(uuid.UUID)
		if key, err = k.murmurShortener.GetUrlId(context.Background(), uuidObj.String()); err != nil {
			logger.LoadExtra(map[string]interface{}{
				"error": err,
			}).Error("GenerateKey: error")
		}
		uuidPool.Put(uuidObj)
		isExist := k.keyData.CheckKeyExist(key)
		if isExist == 1 {
			continue
		}
		k.storeBuffer <- key

	}
}

func (k *KeyGenerator) storeKey() {
	keys := []models.KeyRow{}
	t := time.NewTimer(10 * time.Second)
Loop:
	for {
		select {
		case key, ok := <-k.storeBuffer:
			if !ok {
				if _, err := k.keyData.InsertAvailableKey(keys); err != nil {
					logger.LoadExtra(map[string]interface{}{
						"error": err,
					}).Error("StoreKey: error")
				}
				break Loop
			}

			if len(keys) == k.storeBatchSize {
				if _, err := k.keyData.InsertAvailableKey(keys); err != nil {
					logger.LoadExtra(map[string]interface{}{
						"error": err,
					}).Error("StoreKey: error")
				}
				keys = []models.KeyRow{}
			}

			keys = append(keys, models.KeyRow{
				Key: key,
			})

		case <-t.C:
			if _, err := k.keyData.InsertAvailableKey(keys); err != nil {
				logger.LoadExtra(map[string]interface{}{
					"error": err,
				}).Error("StoreKey: error")
			}
			keys = []models.KeyRow{}
		}
	}
}

func (k *KeyGenerator) Close() error {
	k.stopChan <- struct{}{}
	close(k.storeBuffer)
	return k.keyData.Close()
}
func max(a []int64) int64 {
	var max int64
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}
