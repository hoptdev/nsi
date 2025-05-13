package psql

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	log    *slog.Logger
	dbPool *pgxpool.Pool
}

func New(log *slog.Logger, connect string) (*Storage, error) {
	config, err := pgxpool.ParseConfig(connect)

	if err != nil {
		return nil, err
	}

	config.MaxConns = 50
	config.MinConns = 10
	config.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Storage{
		log:    log,
		dbPool: dbPool,
	}, nil
}
