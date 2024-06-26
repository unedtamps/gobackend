package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
	"github.com/unedtamps/gobackend/pkg/api"
)

func main() {
	if os.Getenv("ENV") != "prod" {
		if err := godotenv.Load(".env"); err != nil {
			panic(err)
		}
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	ctx := context.Background()
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := api.NewServer(db)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
