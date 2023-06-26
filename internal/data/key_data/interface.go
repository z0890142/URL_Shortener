package key_data

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/models"
)

type KeyData interface {
	InsertAvailableKey(keyRows []models.KeyRow) (int, error)
	GetAvailableKey(num int) ([]models.KeyRow, error)
	DeleteAvailableKey(keys []string) error
	InsertAllocatedKey(keyRows []models.KeyRow) (int, error)
	GetAvailableKeyCount() (int64, error)
	CheckKeyExist(key string) int64
	Close() error
}

func NewKeyData(conf config.DatabaseOption) (KeyData, error) {

	return newDefaultKeyData(conf)
}
