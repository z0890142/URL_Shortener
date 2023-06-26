package hook

import (
	"URL_Shortener/pkg/app"
	"URL_Shortener/pkg/utils/common"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormMysql "gorm.io/driver/mysql"

	"gorm.io/gorm"
)

func InitDatabaseHook(app *app.Application) error {
	db, err := common.OpenMysqlDatabase(&app.GetConfig().Databases)
	if err != nil {
		return fmt.Errorf("InitDatabaseHook: %s", err)
	}
	app.SetDatabase(db)
	if err := migration(app); err != nil {
		return fmt.Errorf("InitDatabaseHook: %s", err)
	}
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

func migration(app *app.Application) error {
	driver, err := mysql.WithInstance(app.GetDatabase(), &mysql.Config{})
	if err != nil {
		return fmt.Errorf("Migrate: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", app.GetConfig().MigrationFilePath),
		app.GetConfig().Databases.DBName,
		driver)
	if err != nil {
		return fmt.Errorf("Migrate: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Migrate: %v", err)
	}

	app.GetLogger().Info("Migrate: Migrate successfully")
	return nil
}
