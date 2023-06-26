package key_data

import (
	"URL_Shortener/c"
	"URL_Shortener/config"
	"URL_Shortener/internal/models"
	"URL_Shortener/pkg/utils/common"
	"fmt"

	"github.com/AuthMe01/authme-go-kit/logger"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm/clause"

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

	return &defaultKeyData{
		gormClient: gormClient,
	}, nil
}

func (e *defaultKeyData) InsertAvailableKey(keyRows []models.KeyRow) (int, error) {
	tx := e.gormClient.Begin()
	err := tx.Clauses(clause.Insert{Modifier: "IGNORE"}).
		Table(c.AvailableKey).Create(&keyRows).Error
	if err != nil {
		return 0, fmt.Errorf("InsertKey: %w", tx.Error)
	}

	err = tx.Commit().Error
	if err != nil {
		return 0, fmt.Errorf("InsertKey: %w", tx.Error)
	}

	logger.LoadExtra(map[string]interface{}{
		"inserted": tx.RowsAffected,
	}).Info("InsertKey: Inserted")
	return int(tx.RowsAffected), nil
}

func (e *defaultKeyData) GetAvailableKey(num int) ([]models.KeyRow, error) {
	keyRows := make([]models.KeyRow, 0)
	err := e.gormClient.
		Table(c.AvailableKey).
		Limit(num).Find(&keyRows).Error
	if err != nil {
		return nil, fmt.Errorf("GetKey: %w", err)
	}
	return keyRows, nil
}

func (e *defaultKeyData) DeleteAvailableKey(keys []string) error {
	err := e.gormClient.
		Table(c.AvailableKey).
		Where("`key` in ?", keys).
		Delete(&models.KeyRow{}).Error
	if err != nil {
		return fmt.Errorf("DeleteKey: %w", err)
	}
	return nil
}

func (e *defaultKeyData) InsertAllocatedKey(keyRows []models.KeyRow) (int, error) {
	tx := e.gormClient.
		Clauses(clause.Insert{Modifier: "IGNORE"}).
		Table(c.AllocatedKey).Create(&keyRows)
	if tx.Error != nil {
		return 0, fmt.Errorf("UpdateKey: %w", tx.Error)
	}
	return int(tx.RowsAffected), nil
}

func (e *defaultKeyData) GetAvailableKeyCount() (int64, error) {
	var count int64
	err := e.gormClient.
		Table(c.AvailableKey).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("GetKeyCount: %w", err)
	}
	return count, nil
}

func (e *defaultKeyData) CheckKeyExist(key string) int64 {
	var count int64
	e.gormClient.Table(c.AllocatedKey).Where("`key` = ?", key).Count(&count)
	if count > 0 {
		return count
	}
	e.gormClient.Table(c.AvailableKey).Where("`key` = ?", key).Count(&count)
	return count
}
func (e *defaultKeyData) Close() error {
	db, err := e.gormClient.DB()
	if err != nil {
		return fmt.Errorf("Close: %w", err)
	}
	return db.Close()
}
