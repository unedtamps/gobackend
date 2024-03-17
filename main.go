package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/unedtamps/gobackend/pkg/api"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	fmt.Println("Hello")
	fmt.Println(os.Getenv("POSTGRES_HOST"))

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
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
