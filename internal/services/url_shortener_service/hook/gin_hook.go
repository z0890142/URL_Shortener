package hook

import (
	"URL_Shortener/internal/services/url_shortener_service/controller"
	"URL_Shortener/internal/services/url_shortener_service/handler"
	"URL_Shortener/pkg/app"
	"URL_Shortener/pkg/middleware"
	"fmt"

	"URL_Shortener/pkg/utils/trace"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var defaultController *controller.ShortenerController

func initCtrl(app *app.Application, r *gin.Engine) (*controller.ShortenerController, error) {

	handlerConf := handler.DefaultHandlerConf{
		HashPoolSize: app.GetConfig().HashPoolSize,
		RedisOpts:    app.GetConfig().Redis,
		DatabaseOpts: app.GetConfig().Databases,
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
	if app.GetConfig().Trace.Enable {
		trace.NewTracerProvider(app.GetConfig().Trace.Endpoint, "url_shortener")
		r.Use(otelgin.Middleware("url_shortener"))
	}

	if app.GetConfig().Ratelimit.Enable {
		r.Use(middleware.Ratelimit(middleware.RatelimitConfig{
			Second: app.GetConfig().Ratelimit.Secend,
			Number: app.GetConfig().Ratelimit.Number,
		}))
	}

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
