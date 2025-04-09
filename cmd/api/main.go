package main

import (
	"context"
	"net/http"
	"os"

	"github.com/bernardinorafael/go-boilerplate/internal/infra/database/pg"
	"github.com/bernardinorafael/go-boilerplate/internal/infra/http/middleware"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/auth"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/session"
	"github.com/bernardinorafael/go-boilerplate/internal/modules/user"
	"github.com/bernardinorafael/go-boilerplate/pkg/config"
	"github.com/bernardinorafael/go-boilerplate/pkg/logging"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
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
	r.Use(chimiddleware.Logger)
	r.Use(middleware.WithIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	con, err := pg.NewConnection(cfg.PostgresDSN)
	if err != nil {
		log.Criticalw(ctx, "failed to connect database", logging.Err(err))
		panic(err)
	}
	defer con.Close()

	// User
	userRepo := user.NewRepo(con.DB())
	userService := user.NewService(log, userRepo)

	// Session
	sessionRepo := session.NewRepo(con.DB())
	sessionService := session.NewService(log, sessionRepo, userService, cfg.JWTSecretKey)

	// Auth
	authService := auth.NewService(log, userService, sessionService, cfg.JWTSecretKey)
	auth.NewHandler(authService, cfg.JWTSecretKey).Register(r)

	log.Info(ctx, "Server running")
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Criticalw(ctx, "failed to start server", logging.Err(err))
		os.Exit(1)
	}
}
