package main

import (
	"context"
	"go-templ/internal/handler"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go-templ/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run todo divide everything into layers
func run() error {
	logger, _ := zap.NewProduction()

	//nolint:errcheck
	defer logger.Sync()

	logger.Info("starting app")
	logger.Info("initializing config")

	cfg, err := config.Load()
	if err != nil {
		logger.Panic("cant init config", zap.Error(err))
	}

	logger.Info("initializing database")

	conn, err := sqlx.Open("postgres", cfg.DBUrl)
	if err != nil {
		logger.Panic("cant init db", zap.Error(err))
	}

	if err = conn.Ping(); err != nil {
		logger.Panic("cant connect to db", zap.Error(err))
	}

	logger.Info("initializing router")

	router := mux.NewRouter()
	router.HandleFunc("/", handler.HomeHandler).Methods("GET")

	logger.Info("initializing server")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         net.JoinHostPort("", cfg.ApiPort),
		Handler:      router,
		ReadTimeout:  cfg.HTTPTimeout,
		WriteTimeout: cfg.HTTPTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		logger.Info("starting server", zap.String("port", cfg.ApiPort))

		if err = srv.ListenAndServe(); err != nil {
			logger.Error("failed to start server", zap.Error(err))
		}
	}()

	logger.Info("server started successfully")

	<-done
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	logger.Info("shutting down server", zap.Duration("timeout", cfg.ShutdownTimeout))

	if err = srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("closing database connection")
	if err := conn.Close(); err != nil {
		logger.Error("failed to close db connection", zap.Error(err))
	}

	logger.Info("server stopped")

	return nil
}
