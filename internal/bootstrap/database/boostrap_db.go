package database

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unedtamps/gobackend/internal/config"
	"github.com/unedtamps/gobackend/pkg/utils"
)

type DB struct {
	PRIMARY   *pgxpool.Pool
	SECONDARY *sql.DB
	BACKUP    *sql.DB
}

func ConnectPG(ctx context.Context, conf config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := utils.LoadPostgresConnString(
		conf.Port,
		conf.User,
		conf.Host,
		conf.Password,
		conf.DBName,
	)
	return newPostgressPool(ctx, connStr) // Assuming this function exists in your package
}

func ConnectMySQL(ctx context.Context, conf config.MySQLConfig) (*sql.DB, error) {
	connStr := utils.LoadMySQLConnString(
		conf.Port,
		conf.User,
		conf.Host,
		conf.Password,
		conf.DBName,
	)
	return newMysqlPool(ctx, connStr) // Assuming this function exists in your package
}

func ConnectSQLite(ctx context.Context, conf config.SQLiteConfig) (*sql.DB, error) {
	connStr := utils.LoadSqliteConnString(conf.Path)
	return newSQLitePool(ctx, connStr) // Assuming this function exists in your package
}

func NewDBInstance(ctx context.Context, cfg *config.Config) (*DB, error) {
	var err error
	database := &DB{}
	database.PRIMARY, err = ConnectPG(ctx, cfg.Databases.Primary)
	if err != nil {
		return nil, err
	}

	database.SECONDARY, err = ConnectMySQL(ctx, cfg.Databases.Secondary)
	if err != nil {
		return nil, err
	}
	database.BACKUP, err = ConnectSQLite(ctx, cfg.Databases.Backup)
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (db *DB) Close() {
	if db == nil {
		return
	}
	db.PRIMARY.Close()
	db.SECONDARY.Close()
	db.BACKUP.Close()
}
