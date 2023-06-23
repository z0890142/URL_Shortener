package app

import (
	"URL_Shortener/pkg/utils/logger"
)

func initLoggerApplicationHook(app *Application) error {
	l, err := logger.New(logger.Options{
		Level:   app.GetConfig().LogLevel,
		Outputs: []string{app.GetConfig().LogFile},
	})

	if err != nil {
		panic(err)
	}
	logger.SetLogger(l)
	return nil
}
