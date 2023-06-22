package models

import "time"

type UrlMappingRow struct {
	UrlId       string    `gorm:"column:url_id"`
	OriginalUrl string    `gorm:"column:original_url"`
	Expired     int       `gorm:"column:expired"`
	ExpiredAt   time.Time `gorm:"column:expired_at"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type KeyRow struct {
	Id   int64  `gorm:"column:id"`
	Key  string `gorm:"column:key"`
	Used int    `gorm:"column:used"`
}
