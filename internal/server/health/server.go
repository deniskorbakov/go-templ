package health

import (
	"context"
	"net"
	"net/http"

	"github.com/jmoiron/sqlx"
	"go-templ/internal/resource"
	"go.uber.org/zap"
)

type Server struct {
	server   http.Server
	postgres *sqlx.DB
	logger   *zap.Logger
	errors   chan error
}

func NewServer(cfg *resource.Config, log *zap.Logger, pg *sqlx.DB) *Server {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ready", readyHandler(log, pg))

	return &Server{
		server: http.Server{
			Addr:    net.JoinHostPort("", cfg.DiagPort),
			Handler: nil,
		},
		postgres: pg,
		errors:   make(chan error, 1),
		logger:   log,
	}
}

func (s *Server) Notify() <-chan error {
	return s.errors
}

func (s *Server) Start() {
	go func() {
		s.errors <- s.server.ListenAndServe()
		close(s.errors)
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readyHandler(logger *zap.Logger, pg *sqlx.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := pg.Ping(); err != nil {
			logger.Error("cant pinging postgres", zap.Error(err))
		}

		w.WriteHeader(http.StatusOK)
	}
}
