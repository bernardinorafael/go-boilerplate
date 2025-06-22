package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
)

type Config struct {
	Log          *log.Logger
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Router       *chi.Mux
}

type Server struct {
	log    *log.Logger
	server *http.Server
	config Config
}

func New(c Config) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", c.Port),
		IdleTimeout:  c.IdleTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		Handler:      c.Router,
	}

	return &Server{
		log:    c.Log,
		server: srv,
		config: c,
	}
}

func (s *Server) Start() error {
	s.log.Info("server started", "port", s.config.Port)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
