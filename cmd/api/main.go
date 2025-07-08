package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/config"
	"github.com/bernardinorafael/go-boilerplate/internal/domain/code"
	"github.com/bernardinorafael/go-boilerplate/internal/domain/product"
	"github.com/bernardinorafael/go-boilerplate/internal/domain/user"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/pg"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/mail"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"
	"github.com/bernardinorafael/go-boilerplate/pkg/server"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	env := config.GetConfig()
	metrics := metric.New()
	ctx := context.Background()

	logger := log.NewWithOptions(
		os.Stdout,
		log.Options{
			// Prefix:          env.AppName,
			TimeFormat:      time.Kitchen,
			Formatter:       log.JSONFormatter,
			ReportTimestamp: true,
		},
	)

	if env.Debug {
		logger.SetLevel(log.DebugLevel)
		logger.SetReportCaller(true)
	}
	// in development environment we use TextFormatter
	// to make it easier to read logs in the console
	if env.Environment == "development" {
		logger.SetFormatter(log.TextFormatter)
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic",
				"panic", r,
				"stack", string(debug.Stack()),
			)
		}
	}()

	r := chi.NewMux()
	middleware.Apply(r, middleware.Config{Metrics: metrics})
	r.Handle("/metrics", promhttp.HandlerFor(metrics.Registry(), promhttp.HandlerOpts{}))

	cache, err := cache.New(ctx, env.RedisHost, env.RedisPort, env.RedisPassword)
	if err != nil {
		logger.Error("failed to connect to cache", "error", err)
		panic(err)
	}
	defer cache.Close()

	pgconn, err := pg.NewConnection(env.PostgresDSN)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		panic(err)
	}
	defer pgconn.Close()

	// Repositories
	timeout := time.Second * 2

	prodRepo := product.NewRepo(pgconn.DB(), timeout)
	userRepo := user.NewRepo(pgconn.DB(), timeout)
	codeRepo := code.NewRepo(pgconn.DB(), timeout)

	// Services
	mailService := mail.New(ctx, logger, env.ResendKey, time.Second*5)
	codeService := code.NewService(
		code.ServiceConfig{
			Log:      logger,
			CodeRepo: codeRepo,
			Metrics:  metrics,
			Cache:    cache,
			Mail:     mailService,
		},
	)
	prodService := product.NewService(
		product.ServiceConfig{
			Log:         logger,
			ProductRepo: prodRepo,
			Metrics:     metrics,
			Cache:       cache,
		},
	)
	userService := user.NewService(
		user.ServiceConfig{
			Log:     logger,
			Metrics: metrics,
			Cache:   cache,
			Mail:    mailService,

			UserRepo:    userRepo,
			CodeService: codeService,

			AccessTokenDuration: env.JWTAccessTokenDuration,
			SecretKey:           env.JWTSecretKey,
		},
	)

	// Handlers
	product.NewHandler(prodService, env.JWTSecretKey).Register(r)
	user.NewHandler(userService, env.JWTSecretKey).Register(r)

	srv := server.New(server.Config{
		Log:          logger,
		Port:         env.Port,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		Router:       r,
	})

	shutdownErr := srv.GracefulShutdown(ctx, time.Second*30)

	err = srv.Start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("failed to start server", "error", err)
		os.Exit(1)
	}

	err = <-shutdownErr
	if err != nil {
		logger.Fatal("failed to shutdown server gracefully", "error", err)
		os.Exit(1)
	}

	logger.Info("server shutdown gracefully")
}
