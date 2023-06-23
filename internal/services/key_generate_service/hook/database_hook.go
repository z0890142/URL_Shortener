package hook

import (
	"URL_Shortener/pkg/app"
	"URL_Shortener/pkg/utils/common"
	"fmt"

	gormMysql "gorm.io/driver/mysql"

	"gorm.io/gorm"
)

func InitDatabaseHook(app *app.Application) error {
	db, err := common.OpenMysqlDatabase(&app.GetConfig().Databases)
	if err != nil {
		return fmt.Errorf("InitDatabaseHook: %s", err)
	}
	app.SetDatabase(db)
	return nil
}

func InitGormClientHook(app *app.Application) error {
	gormClient, err := gorm.Open(gormMysql.New(gormMysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      app.GetDatabase(),
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		return fmt.Errorf("InitGormClientHook: %v", err)
	}
	app.SetGormClient(gormClient)
	return nil
}
