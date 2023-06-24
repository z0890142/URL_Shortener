package hook

import (
	"URL_Shortener/internal/services/url_shortener_service/controller"
	"URL_Shortener/internal/services/url_shortener_service/handler"
	"URL_Shortener/pkg/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

var defaultController *controller.ShortenerController

func initCtrl(app *app.Application, r *gin.Engine) (*controller.ShortenerController, error) {

	handlerConf := handler.DefaultHandlerConf{
		HashPoolSize: app.GetConfig().HashPoolSize,
		RedisOpts:    app.GetConfig().Redis,
		DatabaseOpts: app.GetConfig().Databases,
		RetryTimes:   app.GetConfig().RetryTimes,
	}

	if app.GetConfig().EnableKeyService {
		keyServerAdd := fmt.Sprintf("%s:%s",
			app.GetConfig().EndPoints.KeyServer.Http.Host,
			app.GetConfig().EndPoints.KeyServer.Http.Port)

		if app.GetConfig().EndPoints.KeyServer.Http.EnableTls {
			keyServerAdd = "https://" + keyServerAdd
		} else {
			keyServerAdd = "http://" + keyServerAdd
		}
		handlerConf.KeyServiceAddr = keyServerAdd
		handlerConf.EnableKeyService = true
	}

	handler, err := handler.NewDefaultShortenerHandler(handlerConf)

	if err != nil {
		panic(err)
	}

	defaultController = controller.NewController(handler)

	r.POST("/api/v1/urls", defaultController.NewShortUrl)
	r.GET("/:urlId", defaultController.RedirectUrl)

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
	addr := fmt.Sprintf("%s:%s", app.GetConfig().ShortenerService.Host, app.GetConfig().ShortenerService.Port)

	app.SetAddr(addr)
	app.SetSrv(r)

	return nil
}

func DestroyGinApplicationHook(app *app.Application) error {
	defaultController.Shutdown()
	return nil
}
