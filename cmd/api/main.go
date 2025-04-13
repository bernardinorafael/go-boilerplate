package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/config"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/pg"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/server"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/mail"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/auth"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/product"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/session"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/user"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/metric"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.GetConfig()
	metrics := metric.New()
	ctx := context.Background()

	r := chi.NewRouter()
	middleware.Apply(r, middleware.Config{
		Metrics: metrics,
	})
	r.Handle("/metrics", promhttp.HandlerFor(metrics.Registry(), promhttp.HandlerOpts{}))

	cache, err := cache.New(ctx, cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	if err != nil {
		slog.Error("failed to connect to cache", "error", err)
		panic(err)
	}
	defer cache.Close()

	pgconn, err := pg.NewConnection(cfg.PostgresDSN)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		panic(err)
	}
	defer pgconn.Close()

	// Repositories
	userRepo := user.NewRepo(pgconn.DB())
	sessRepo := session.NewRepo(pgconn.DB())
	prodRepo := product.NewRepo(pgconn.DB())

	// Services
	mailService := mail.New(ctx, mail.Config{
		APIKey:     cfg.ResendKey,
		RetryDelay: time.Second * 2,
		Timeout:    time.Second * 5,
		MaxRetries: 3,
	})
	userService := user.NewService(user.ServiceConfig{
		UserRepo: userRepo,
	})
	sessionService := session.NewService(session.ServiceConfig{
		SessionRepo: sessRepo,
		UserService: userService,
		Cache:       cache,
		Metrics:     metrics,
		SecretKey:   cfg.JWTSecretKey,
	})
	prodService := product.NewService(product.ServiceConfig{
		ProductRepo: prodRepo,
		Metrics:     metrics,
		Cache:       cache,
	})
	authService := auth.NewService(auth.ServiceConfig{
		UserService:    userService,
		UserRepo:       userRepo,
		SessionService: sessionService,
		SessionRepo:    sessRepo,
		Mailer:         mailService,
		Cache:          cache,
		Metrics:        metrics,
		SecretKey:      cfg.JWTSecretKey,
	})

	// Handlers
	auth.NewHandler(authService, cfg.JWTSecretKey).Register(r)
	session.NewHandler(sessionService, cfg.JWTSecretKey).Register(r)
	product.NewHandler(prodService, cfg.JWTSecretKey).Register(r)

	srv := server.New(server.Config{
		Port:         cfg.Port,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		Router:       r,
	})

	shutdoewnErr := srv.GracefulShutdown(ctx, time.Second*30)

	err = srv.Start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	err = <-shutdoewnErr
	if err != nil {
		slog.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	}

	slog.Info("server shutdown gracefully")
}
