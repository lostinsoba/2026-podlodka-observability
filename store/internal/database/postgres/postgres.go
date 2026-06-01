package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"external/sdk/collector"
	"github.com/prometheus/client_golang/prometheus"
	"store/internal/config"
)

type Database struct {
	db *sql.DB
}

func New(cfg config.Database, mr prometheus.Registerer) (*Database, error) {
	db, err := sql.Open("postgres", cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	err = collector.RegisterDBMetricCollector(mr, db)
	if err != nil {
		return nil, fmt.Errorf("failed to register db metric collector: %w", err)
	}
	return &Database{db: db}, nil
}
