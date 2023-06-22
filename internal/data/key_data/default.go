package key_data

import (
	"URL_Shortener/c"
	"URL_Shortener/config"
	"URL_Shortener/internal/models"
	"URL_Shortener/internal/utils/common"
	"fmt"

	gormMysql "gorm.io/driver/mysql"

	"gorm.io/gorm"
)

type defaultKeyData struct {
	gormClient *gorm.DB
}

func newDefaultKeyData(conf config.DatabaseOption) (KeyData, error) {
	db, err := common.OpenMysqlDatabase(&conf)
	if err != nil {
		return nil, fmt.Errorf("NewUrlMappingMysql: %s", err)
	}
	gormClient, err := gorm.Open(gormMysql.New(gormMysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		return nil, fmt.Errorf("NewGormClient: %v", err)
	}

	return &defaultKeyData{gormClient: gormClient}, nil
}

func (e *defaultKeyData) InsertKey(keyRows []models.KeyRow) (int, error) {
	tx := e.gormClient.
		Table(c.Key).Create(&keyRows)
	if tx.Error != nil {
		return 0, fmt.Errorf("InsertKey: %w", tx.Error)
	}
	return int(tx.RowsAffected), nil
}

func (e *defaultKeyData) GetKey(num int, startId int64) ([]models.KeyRow, error) {
	keyRows := make([]models.KeyRow, 0)
	err := e.gormClient.
		Table(c.Key).Where("id >= ?", startId).Limit(num).Find(&keyRows).Error
	if err != nil {
		return nil, fmt.Errorf("GetKey: %w", err)
	}
	return keyRows, nil
}

func (e *defaultKeyData) UpdateKey(keyRows []models.KeyRow) (int, error) {
	tx := e.gormClient.
		Table(c.Key).Save(&keyRows)
	if tx.Error != nil {
		return 0, fmt.Errorf("UpdateKey: %w", tx.Error)
	}
	return int(tx.RowsAffected), nil
}
