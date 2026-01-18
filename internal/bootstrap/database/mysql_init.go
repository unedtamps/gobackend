package database

import (
	"context"
	"database/sql"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func newMysqlPool(ctx context.Context, connString string) (*sql.DB, error) {
	mysqlPool, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	mysqlPool.SetMaxOpenConns(runtime.NumCPU() * 4)
	mysqlPool.SetMaxIdleConns(runtime.NumCPU() * 2)

	mysqlPool.SetConnMaxIdleTime(3 * time.Minute)
	mysqlPool.SetConnMaxLifetime(30 * time.Minute)
	err = mysqlPool.PingContext(ctx)
	return mysqlPool, err
}
