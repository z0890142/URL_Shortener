package key_data

import "URL_Shortener/internal/models"

type KeyData interface {
	GetKey(num int, startId int64) ([]models.KeyRow, error)
	InsertKey([]models.KeyRow) (int, error)
	UpdateKey([]models.KeyRow) (int, error)
}
