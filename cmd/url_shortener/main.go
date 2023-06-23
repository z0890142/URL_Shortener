package main

import (
	"URL_Shortener/config"
	"URL_Shortener/internal/services/url_shortener_service/hook"
	"URL_Shortener/pkg/app"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../", "config path, eg: -conf config.yaml")
}

func handleSignals(server *app.Application) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	server.GetLogger().Infof("signal %s received", <-sigs)
	server.Shutdown()
}

func main() {
	flag.Parse()
	config.LoadConf(flagconf, config.GetConfig())

	server := app.Default()
	server.AddInitHook(hook.InitGinApplicationHook)

	go handleSignals(server)
	server.Run()
}
