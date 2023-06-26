package hook

import (
	"URL_Shortener/internal/services/key_generate_service/controller"
	"URL_Shortener/internal/services/key_generate_service/handler"
	"URL_Shortener/pkg/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

var defaultController *controller.KeyController

func initCtrl(app *app.Application, r *gin.Engine) (*controller.KeyController, error) {

	handler, err := handler.NewKeyHandler(handler.DefaultKeyHandlerConf{
		HashPoolSize:   app.GetConfig().HashPoolSize,
		StoreBatchSize: app.GetConfig().StoreBatchSize,
	})

	if err != nil {
		return nil, err
	}
	defaultController = controller.NewController(handler)
	r.POST("/api/v1/key", defaultController.GenerateKey)
	return defaultController, nil
}

func InitGinApplicationHook(app *app.Application) error {
	if app.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.EnableJsonDecoderUseNumber()
	r := gin.New()
	r.Use(gin.Recovery())

	initCtrl(app, r)
	addr := fmt.Sprintf("%s:%s", app.GetConfig().KeyService.Host, app.GetConfig().KeyService.Port)

	app.SetAddr(addr)
	app.SetSrv(r)

	return nil
}

func DestroyGinApplicationHook(app *app.Application) error {
	defaultController.Shutdown()
	return nil
}
