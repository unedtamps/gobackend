package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/unedtamps/gobackend/config"
	server "github.com/unedtamps/gobackend/src"
)

func main() {
	connStr := config.ConStr()
	ctx := context.Background()
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		log.Fatal(err)
	}
	quit := make(chan os.Signal)

	go func(db *pgxpool.Pool) {
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit
		db.Close()
		os.Exit(0)
	}(db)

	s := server.NewServer(db)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
