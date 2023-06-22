package main

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/services/key_generate_service/controller"
	"URL_Shortener/internal/utils/logger"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	config.LoadConf(flagconf, config.GetConfig())
	initLogger()
	initGin()
}

func initLogger() {
	l, err := logger.New(logger.Options{
		Level:   config.GetConfig().LogLevel,
		Outputs: []string{config.GetConfig().LogFile},
	})

	if err != nil {
		panic(err)
	}
	logger.SetLogger(l)

}

func initGin() {
	if config.GetConfig().Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.EnableJsonDecoderUseNumber()
	r := gin.New()

	r.Use(gin.Recovery())
	controller.RegisterRoutes(r)

	addr := fmt.Sprintf("%s:%s", config.GetConfig().KeyService.Host, config.GetConfig().KeyService.Port)
	if err := r.Run(addr); err != nil {
		logger.LoadExtra(map[string]interface{}{
			"addr":  addr,
			"error": err,
		}).Error("run gin http server error")
		panic(err)
	}
	logger.Debug("init gin http server")

}
