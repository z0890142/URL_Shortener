package handler

import (
	"URL_Shortener/internal/data/key_data"
	"URL_Shortener/internal/models"
	"URL_Shortener/internal/services/key_generate_service/key_generator"
	"URL_Shortener/pkg/app"
	"URL_Shortener/pkg/utils/shortener"
	"fmt"
)

type defaultKeyHandler struct {
	keyData      key_data.KeyData
	keyGenerator *key_generator.KeyGenerator
}

type DefaultKeyHandlerConf struct {
	HashPoolSize   int
	StoreBatchSize int
}

func newDefaultKeyHandler(conf DefaultKeyHandlerConf) (KeyHandler, error) {

	defaultKeyHandler := defaultKeyHandler{}

	murmurShortener := shortener.NewShortener(shortener.MurMurShortenerConfig{
		HashPoolSize: conf.HashPoolSize,
	})

	keyData, err := key_data.NewKeyData(app.Default().GetConfig().Databases)
	if err != nil {
		return nil, fmt.Errorf("NewDefaultKeyHandler: %w", err)
	}
	defaultKeyHandler.keyData = keyData

	keyGenerator := key_generator.NewKeyGenerator(conf.StoreBatchSize, keyData, murmurShortener)
	defaultKeyHandler.keyGenerator = keyGenerator

	defaultKeyHandler.keyGenerator.Start()
	return &defaultKeyHandler, nil
}

func (d *defaultKeyHandler) GetKeys(num int) (result []string, err error) {

	var keys []models.KeyRow

	keys, err = d.keyData.GetAvailableKey(num)
	if err != nil {
		return []string{}, fmt.Errorf("GetKeys: %w", err)
	}

	keyIds := make([]string, len(keys))
	for i, key := range keys {
		keyIds[i] = key.Key
	}
	err = d.keyData.DeleteAvailableKey(keyIds)
	if err != nil {
		return []string{}, fmt.Errorf("GetKeys: %w", err)
	}

	_, err = d.keyData.InsertAllocatedKey(keys)
	if err != nil {
		return []string{}, fmt.Errorf("GetKeys: %w", err)
	}

	result = make([]string, len(keys))
	for i, key := range keys {
		result[i] = key.Key
	}

	return result, nil
}

func (d *defaultKeyHandler) Shutdown() {
	d.keyData.Close()
}
