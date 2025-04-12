package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/pg"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/mail"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/auth"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/product"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/session"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/user"
	"github.com/bernardinorafael/go-boilerplate/pkg/cache"
	"github.com/bernardinorafael/go-boilerplate/pkg/config"
	"github.com/bernardinorafael/go-boilerplate/pkg/logging"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfig()

	log := logging.New(logging.LogParams{
		AppName:                  cfg.AppName,
		DebugLevel:               cfg.DebugMode,
		AddAttributesFromContext: nil,
		LogToFile:                false,
	})

	r := chi.NewRouter()
	r.Use(
		// cmid.Logger,
		middleware.WithIP,
		middleware.WithRateLimit,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	cache, err := cache.New(ctx, &cache.Config{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})
	if err != nil {
		log.Criticalw(ctx, "failed to connect to cache", logging.Err(err))
		panic(err)
	}
	defer cache.Close()

	// redisconn, err := rdb.NewConnection(ctx, cfg)
	// if err != nil {
	// 	log.Criticalw(ctx, "failed to connect to redis", logging.Err(err))
	// 	panic(err)
	// }
	// defer redisconn.Close()

	pgconn, err := pg.NewConnection(cfg.PostgresDSN)
	if err != nil {
		log.Criticalw(ctx, "failed to connect database", logging.Err(err))
		panic(err)
	}
	defer pgconn.Close()

	// Mailer
	mailService := mail.New(ctx, log, mail.Config{
		MaxRetries: 3,
		APIKey:     cfg.ResendKey,
		RetryDelay: time.Second * 2,
		Timeout:    time.Second * 5,
	})
	// User
	userRepo := user.NewRepo(pgconn.DB())
	userService := user.NewService(log, userRepo)
	// Session
	sessionRepo := session.NewRepo(pgconn.DB())
	sessionService := session.NewService(log, sessionRepo, userService, cache, cfg.JWTSecretKey)
	session.NewHandler(sessionService, cfg.JWTSecretKey).Register(r)
	// Auth
	// TODO: Consider using option pattern to avoid having so many parameters
	authService := auth.NewService(
		log,
		userService,
		userRepo,
		sessionService,
		sessionRepo,
		mailService,
		cache,
		cfg.JWTSecretKey,
	)
	auth.NewHandler(authService, cfg.JWTSecretKey).Register(r)
	// Products
	productRepo := product.NewRepo(pgconn.DB())
	productService := product.NewService(log, productRepo)
	product.NewHandler(productService, cfg.JWTSecretKey).Register(r)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		Handler:      r,
	}
	log.Info(ctx, "Starting server")

	shutdownErr := make(chan error)
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(
			stop,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		sig := <-stop

		log.Infow(ctx,
			"caught signal...",
			logging.String("signal", sig.String()),
		)

		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		shutdownErr <- server.Shutdown(ctx)
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Criticalw(ctx, "failed to start server", logging.Err(err))
		os.Exit(1)
	}

	err = <-shutdownErr
	if err != nil {
		log.Criticalw(ctx, "failed to shutdown server gracefully", logging.Err(err))
		os.Exit(1)
	}

	log.Info(ctx, "Server shutdown gracefully")
}
