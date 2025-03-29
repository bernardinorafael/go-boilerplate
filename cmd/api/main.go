package main

import (
	"context"
	"gulg/internal/infra/database/pg"
	"gulg/internal/infra/http/middleware"
	"gulg/internal/modules/auth"
	"gulg/internal/modules/session"
	"gulg/internal/modules/user"
	"gulg/pkg/config"
	"gulg/pkg/logging"
	"net/http"
	"os"

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
	userService := user.NewService(userRepo)

	// Session
	sessionRepo := session.NewRepo(con.DB())
	sessionService := session.NewService(log, sessionRepo, userService, cfg.JWTSecretKey)

	// Auth
	authService := auth.NewService(log, userService, sessionService, cfg.JWTSecretKey)
	auth.NewHandler(authService).Register(r)

	log.Info(ctx, "Server running")
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Criticalw(ctx, "failed to start server", logging.Err(err))
		os.Exit(1)
	}
}
