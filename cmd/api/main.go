package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/config"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/container"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/handler"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/logger"
	"github.com/bernardinorafael/go-boilerplate/pkg/server"
	"github.com/go-chi/chi/v5"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	env := config.GetConfig()
	ctx := context.Background()

	logger := logger.NewLogger(env)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic",
				"panic", r,
				"stack", string(debug.Stack()),
			)
		}
	}()

	container, err := container.NewContainer(ctx, env, logger)
	if err != nil {
		logger.Fatal("failed to initialize container", "error", err)
		os.Exit(1)
	}
	defer container.Close()

	r := chi.NewMux()
	middleware.Apply(r, middleware.Config{Metrics: container.Metrics})
	r.Handle("/metrics", promhttp.HandlerFor(container.Metrics.Registry(), promhttp.HandlerOpts{}))

	handler.RegisterHandler(r, container)

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
