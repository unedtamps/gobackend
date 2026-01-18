package database

import (
	"context"
	"log"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newPostgressPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	dbConfig.MaxConns = int32(runtime.NumCPU() * 4)
	dbConfig.MinConns = 2

	dbConfig.MaxConnIdleTime = 5 * time.Minute
	dbConfig.MaxConnLifetime = 30 * time.Minute
	dbConfig.HealthCheckPeriod = 1 * time.Minute

	dbConfig.ConnConfig.ConnectTimeout = 5 * time.Second

	postgres, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = postgres.Ping(ctx)
	return postgres, err
}
