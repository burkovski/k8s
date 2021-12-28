package main

import (
	"context"
	"flag"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

const (
	dataFilename = "/var/data/kubia.txt"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "the address on which server will listen for requests")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	go onShutdown(cancel)

	app := echo.New()
	setupApp(app)

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(startApp(app, addr))
	group.Go(shutdownApp(app, groupCtx))

	if err := group.Wait(); err != nil {
		log.Error().Err(err).Msg("exit reason")
	}
}

func onShutdown(handleShutdown func()) {
	waitForShutdown()
	handleShutdown()
}

func waitForShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-signalChan
}

func setupApp(app *echo.Echo) {
	setupPreferences(app)
	setupMiddlewares(app)
	setupRoutes(app)
}

func setupPreferences(app *echo.Echo) {
	app.HideBanner = true
	app.HidePort = true
}

func setupMiddlewares(app *echo.Echo) {
	app.Use(
		middleware.Recover(),
		middleware.RequestID(),
		middleware.Logger(),
	)
}

func setupRoutes(app *echo.Echo) {
	app.POST("/", storeData)
	app.GET("/", getData)
}

func startApp(app *echo.Echo, addr string) func() error {
	return func() error {
		log.Info().Str("addr", addr).
			Msg("starting kubia server")

		return app.Start(addr)
	}
}

func shutdownApp(app *echo.Echo, groupCtx context.Context) func() error {
	return func() error {
		<-groupCtx.Done()
		log.Info().Str("reason", groupCtx.Err().Error()).
			Msg("shutting down kubia server")

		return app.Shutdown(context.Background())
	}
}
