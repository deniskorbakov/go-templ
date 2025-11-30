package main

import (
	"database/sql"
	"log"

	"go-http-template/internal/config"
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

	cfg, err := config.MustLoad()
	if err != nil {
		logger.Panic("cant init config, err: ", zap.Error(err))
	}

	logger.Info("initializing database")

	conn, err := sql.Open("pgx", cfg.DBUrl)
	if err != nil {
		logger.Panic("cant init db, err: ", zap.Error(err))
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			logger.Error("cant close db, err: ", zap.Error(err))
		}
	}()

	return nil
}
