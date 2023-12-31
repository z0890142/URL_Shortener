package url_mapping_data

import (
	"URL_Shortener/c"
	"URL_Shortener/config"
	"URL_Shortener/internal/models"
	"URL_Shortener/pkg/utils/common"
	"URL_Shortener/pkg/utils/logger"
	"URL_Shortener/pkg/utils/trace"
	"context"
	"fmt"
	"time"

	gormMysql "gorm.io/driver/mysql"

	"gorm.io/gorm"
)

type urlMappingMysql struct {
	gormClient *gorm.DB
}

func newUrlMappingMysql(conf config.DatabaseOption) (UrlMappingData, error) {
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

	return &urlMappingMysql{gormClient: gormClient}, nil
}

func (u *urlMappingMysql) SetUrlId(ctx context.Context, urlId, url, expireAt string) error {
	if config.GetConfig().Trace.Enable {
		c, span := trace.NewSpan(ctx, "http://jaeger:14268/api/traces")
		defer span.End()
		ctx = c
	}

	expire, err := time.Parse(c.TimeFormat, expireAt)
	if err != nil {
		return fmt.Errorf("SetUrlId: %w", err)
	}
	urlMappingRow := models.UrlMappingRow{
		UrlId:       urlId,
		OriginalUrl: url,
		ExpiredAt:   expire,
	}

	if err := u.gormClient.WithContext(ctx).Table(c.UrlMapping).Create(&urlMappingRow).Error; err != nil {
		return err
	}
	return nil
}
func (u *urlMappingMysql) GetUrl(ctx context.Context, urlId string) (string, error) {
	if config.GetConfig().Trace.Enable {
		c, span := trace.NewSpan(ctx, "http://jaeger:14268/api/traces")
		defer span.End()
		ctx = c
	}

	urlMappingRow := models.UrlMappingRow{
		UrlId:   urlId,
		Expired: 0,
	}

	err := u.gormClient.Table(c.UrlMapping).Where(&urlMappingRow).First(&urlMappingRow).Error
	if err != nil {
		return "", fmt.Errorf("GetUrl: %w", err)
	}

	if urlMappingRow.ExpiredAt.Before(time.Now()) {
		urlMappingRow.Expired = 1
		err := u.gormClient.Table(c.UrlMapping).Where(&urlMappingRow).Save(&urlMappingRow).Error
		if err != nil {
			logger.LoadExtra(map[string]interface{}{
				"error": err,
				"urlId": urlId,
			}).Error("GetUrl: update expired failed")
		}
		return "", fmt.Errorf("GetUrl: url expired")
	}

	return urlMappingRow.OriginalUrl, nil
}

func (u *urlMappingMysql) Close() error {
	db, err := u.gormClient.DB()
	if err != nil {
		return fmt.Errorf("Close: %w", err)
	}
	return db.Close()
}
