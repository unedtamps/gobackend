package database

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func newSQLitePool(ctx context.Context, connStr string) (*sql.DB, error) {
	return sql.Open("sqlite3", connStr)
}
