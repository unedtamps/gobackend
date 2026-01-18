package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unedtamps/gobackend/internal/config"
	"github.com/unedtamps/gobackend/pkg/utils"
)

type DB struct {
	Pg    map[string]*pgxpool.Pool
	Mysql map[string]*sql.DB
}

func ConnectAll[C any, D any](
	ctx context.Context,
	configs map[string]C,
	connectFn func(context.Context, C) (D, error),
) (map[string]D, error) {
	connections := make(map[string]D)

	for name, cfg := range configs {
		conn, err := connectFn(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s: %w", name, err)
		}
		connections[name] = conn
	}
	return connections, nil
}

func ConnectPG(ctx context.Context, conf config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := utils.LoadPostgresConnString(conf.Port, conf.User, conf.Host, conf.Password, conf.DB)
	return newPostgressPool(ctx, connStr) // Assuming this function exists in your package
}

func ConnectMySQL(ctx context.Context, conf config.MysqlConfig) (*sql.DB, error) {
	connStr := utils.LoadMySQLConnString(conf.Port, conf.User, conf.Host, conf.Password, conf.DB)
	return newMysqlPool(ctx, connStr) // Assuming this function exists in your package
}

func NewDBInstance(ctx context.Context, cfg *config.Config) (*DB, error) {
	var err error
	database := &DB{}

	database.Pg, err = ConnectAll(ctx, cfg.Databases.Postgres, ConnectPG)
	if err != nil {
		return nil, err
	}

	database.Mysql, err = ConnectAll(ctx, cfg.Databases.Mysql, ConnectMySQL)
	if err != nil {
		return nil, err
	}

	if _, ok := database.Pg["primary"]; !ok {
		return nil, fmt.Errorf("postgres 'primary' connection is missing in config")
	}

	if _, ok := database.Mysql["primary"]; !ok {
		return nil, fmt.Errorf("mysql 'primary' connection is missing in config")
	}

	return database, nil
}

func (db *DB) PrimaryPG() *pgxpool.Pool {
	return db.Pg["primary"]
}

func (db *DB) PrimaryMySQL() *sql.DB {
	return db.Mysql["primary"]
}

func (db *DB) Close() {
	if db == nil {
		return
	}
	for _, p := range db.Pg {
		p.Close()
	}
	for _, m := range db.Mysql {
		_ = m.Close()
	}
}
