package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-list-templ/grpc/internal/resource"
	"github.com/go-list-templ/grpc/internal/server/grpc"
	"github.com/go-list-templ/grpc/internal/server/health"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger, _ := zap.NewProduction()

	//nolint:errcheck
	defer logger.Sync()

	logger.Info("starting app")
	logger.Info("initializing config")

	cfg, err := resource.NewConfig()
	if err != nil {
		logger.Panic("cant init config", zap.Error(err))
	}

	logger.Info("initializing postgres")

	pg, err := resource.NewPostgres(cfg)
	if err != nil {
		logger.Panic("cant init postgres", zap.Error(err))
	}

	logger.Info("initializing servers")

	grpcServer := grpc.NewServer(cfg, logger)
	grpcServer.Start()

	healthServer := health.NewServer(cfg, logger, pg)
	healthServer.Start()

	logger.Info("server started successfully")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		logger.Info("Received a signal.", zap.String("signal", x.String()))
	case err = <-healthServer.Notify():
		logger.Error("Received an error from the health server", zap.Error(err))
	case err = <-grpcServer.Notify():
		logger.Error("Received an error from the grpc server", zap.Error(err))
	}

	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	grpcServer.Stop()

	if err = healthServer.Stop(ctx); err != nil {
		logger.Error("server stopped with error", zap.Error(err))
	}

	logger.Info("closing resources")

	if err = pg.Close(); err != nil {
		logger.Error("postgres close failed", zap.Error(err))
	}

	logger.Info("The app is calling the last defers and will be stopped")

	return nil
}
