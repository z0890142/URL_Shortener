package handler

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/data/key_data"
	"URL_Shortener/internal/models"
	"URL_Shortener/internal/utils/logger"
	"URL_Shortener/internal/utils/shortener"
	"fmt"
	"sync"
	"time"
)

type defaultKeyHandler struct {
	mu              *sync.Mutex
	keyBuffer       chan string
	storeBuffer     chan string
	storeBatchSize  int
	latestKeyId     int64
	keyData         key_data.KeyData
	murmurShortener shortener.Shortener
}

type DefaultKeyHandlerConf struct {
	HashPoolSize   int
	StoreBatchSize int
}

func newDefaultKeyHandler(conf DefaultKeyHandlerConf) (KeyHandler, error) {

	defaultKeyHandler := defaultKeyHandler{
		mu:             &sync.Mutex{},
		keyBuffer:      make(chan string, conf.StoreBatchSize),
		storeBuffer:    make(chan string),
		storeBatchSize: conf.StoreBatchSize,
	}

	murmurShortener := shortener.NewShortener(shortener.MurMurShortenerConfig{
		HashPoolSize: conf.HashPoolSize,
	})
	defaultKeyHandler.murmurShortener = murmurShortener

	keyData, err := key_data.NewKeyData(config.GetConfig().Databases)
	if err != nil {
		return nil, fmt.Errorf("NewDefaultKeyHandler: %w", err)
	}
	defaultKeyHandler.keyData = keyData
	defaultKeyHandler.generateKey()

	return &defaultKeyHandler, nil
}

func (d *defaultKeyHandler) GetKeys(num int) (result []string, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var count int
	var keys []models.KeyRow

Loop:
	for count = 0; count < num; count++ {
		select {
		case key := <-d.keyBuffer:
			result = append(result, key)
		default:
			keys, err = d.keyData.GetKey(num*2-count, d.latestKeyId)
			if err != nil {
				d.insertKeyToBuf(result)
				return []string{}, fmt.Errorf("GetKeys: %w", err)
			}
			break Loop
		}
	}

	if len(keys) == 0 {
		return result, nil
	}

	for index, key := range keys {
		keys[index].Used = 1
		if len(result) < num {
			result = append(result, key.Key)
			continue
		}
		d.keyBuffer <- key.Key
	}

	//update keys to used
	if _, err = d.keyData.UpdateKey(keys); err != nil {
		return []string{}, fmt.Errorf("GetKeys: %w", err)
	}

	d.latestKeyId = keys[len(keys)-1].Id

	return result, nil
}

func (d *defaultKeyHandler) generateKey() {
	go d.storeKey()
	for {
		if key, err := d.murmurShortener.GenerateUrlId(time.Now().Format(time.RFC3339Nano)); err != nil {
			logger.LoadExtra(map[string]interface{}{
				"error": err,
			}).Error("GenerateKey: error")
		} else {
			d.storeBuffer <- key
		}
	}
}

func (d *defaultKeyHandler) storeKey() {
	keys := []models.KeyRow{}
Loop:
	for {
		select {
		case key, ok := <-d.storeBuffer:
			if !ok {
				if _, err := d.keyData.InsertKey(keys); err != nil {
					logger.LoadExtra(map[string]interface{}{
						"error": err,
					}).Error("StoreKey: error")
				}
				break Loop
			}

			if len(keys) == d.storeBatchSize {
				if _, err := d.keyData.InsertKey(keys); err != nil {
					logger.LoadExtra(map[string]interface{}{
						"error": err,
					}).Error("StoreKey: error")
				}
				keys = []models.KeyRow{}
			}

			keys = append(keys, models.KeyRow{
				Key:  key,
				Used: 0,
			})

		case <-time.After(30 * time.Second):
			if _, err := d.keyData.InsertKey(keys); err != nil {
				logger.LoadExtra(map[string]interface{}{
					"error": err,
				}).Error("StoreKey: error")
			}
			keys = []models.KeyRow{}
		}
	}
}

func (d *defaultKeyHandler) insertKeyToBuf(keys []string) {
	for _, key := range keys {
		d.keyBuffer <- key
	}
}
