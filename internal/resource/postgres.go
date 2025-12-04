package resource

import (
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

func NewPostgres(cfg *Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", cfg.DBUrl)
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(cfg.DbMaxIdle)
	conn.SetMaxOpenConns(cfg.DbMaxConn)
	conn.SetConnMaxLifetime(cfg.DbMaxConnTime)
	conn.SetConnMaxIdleTime(cfg.DbMaxIdleTime)

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
