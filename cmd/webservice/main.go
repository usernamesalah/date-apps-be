package main

import (
	"context"
	"date-apps-be/infrastructure/config"
	"date-apps-be/infrastructure/database"
	"date-apps-be/internal/api/http/router"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"time"

	"date-apps-be/internal/container"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {

	// Load configuration
	config.Init()

	// Logging
	configLog := zap.NewProductionConfig()
	configLog.EncoderConfig.StacktraceKey = "" // to hide stacktrace info
	configLog.DisableCaller = true

	log, err := configLog.Build()
	if err != nil {
		panic(err)
	}

	if err := run(log); err != nil {
		log.Error("error: shutting down: %s", zap.Error(err))
		os.Exit(1)
	}
}

func run(log *zap.Logger) error {
	conf := config.Get()

	// Start Database
	database, err := database.InitializeDatabase(conf)
	if err != nil {
		log.Error("web failed to init db", zap.Error(err))
		return err
	}

	sharedComponent := &container.SharedComponent{
		DB:   database,
		Conf: conf,
		Log:  log,
	}

	cc := container.NewHandlerComponent(sharedComponent)

	log.Info("Initializing the web server ...")
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Header.Set("Cache-Control", "max-age:3600, public")
			return next(c)
		}
	})

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info("request",
				zap.String("Latency", v.Latency.String()),
				zap.String("Remote IP", c.RealIP()),
				zap.String("URI", v.URI),
				zap.String("Method", c.Request().Method),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	e.Validator = &requestValidator{}

	// init route
	router.Init(e, cc, sharedComponent)

	// Start server
	server := &http.Server{
		Addr:         "0.0.0.0:" + conf.HttpPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	serverErrors := make(chan error, 1)
	// mulai listening server
	go func() {
		log.Info("server listening on", zap.String("address", server.Addr))
		serverErrors <- e.StartServer(server)
	}()

	// Membuat channel untuk mendengarkan sinyal interupsi/terminate dari OS.
	// Menggunakan channel buffered karena paket signal membutuhkannya.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Mengontrol penerimaan data dari channel,
	// jika ada error saat listenAndServe server maupun ada sinyal shutdown yang diterima
	select {
	case err := <-serverErrors:
		return fmt.Errorf("starting server: %v", err)

	case <-shutdown:
		log.Info("caught signal, shutting down")

		// Jika ada shutdown, meminta tambahan waktu 10 detik untuk menyelesaikan proses yang sedang berjalan.
		const timeout = 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Error("error: gracefully shutting down server: %s", zap.Error(err))
			if err := server.Close(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %v", err)
			}
		}

	}

	return nil
}

type requestValidator struct{}

func (rv *requestValidator) Validate(i interface{}) (err error) {
	_, err = govalidator.ValidateStruct(i)
	return
}
