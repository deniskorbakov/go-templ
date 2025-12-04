package health

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go-templ/internal/resource"
	"go.uber.org/zap"
)

type Server struct {
	server   http.Server
	postgres *sqlx.DB
	redis    *redis.Client
	errors   chan error
	logger   *zap.Logger
}

func NewServer(cfg *resource.Config, log *zap.Logger, pg *sqlx.DB, rd *redis.Client) *Server {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ready", readyHandler(log, pg, rd))

	return &Server{
		server: http.Server{
			Addr:    net.JoinHostPort("", cfg.DiagPort),
			Handler: nil,
		},
		postgres: pg,
		redis:    rd,
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

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readyHandler(logger *zap.Logger, pg *sqlx.DB, rd *redis.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := pg.Ping(); err != nil {
			logger.Error("cant pinging postgres", zap.Error(err))
		}

		if err := rd.Ping(ctx).Err(); err != nil {
			logger.Error("cant pinging redis", zap.Error(err))
		}

		w.WriteHeader(http.StatusOK)
	}
}
