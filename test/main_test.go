package test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unedtamps/gobackend/config"
	server "github.com/unedtamps/gobackend/src"
	"github.com/unedtamps/gobackend/src/repository"
)

var repo *repository.Store

func TestMain(m *testing.M) {
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
	defer db.Close()

	go func() {
		s := server.NewServer(db)
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	readyServer := make(chan struct{})

	go func() {
		for {
			resp, err := http.Get(fmt.Sprintf("http://%s:%s",
				config.Config.SERVER_HOST,
				config.Config.SERVER_PORT,
			))
			if err == nil && resp.StatusCode == 200 {
				close(readyServer)
				return
			}
			time.Sleep(time.Millisecond * 100)

		}
	}()
	<-readyServer
	repo = repository.NewStore(db)
	code := m.Run()
	os.Exit(code)
}
