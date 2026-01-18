package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/unedtamps/gobackend/internal/bootstrap/database"
	"github.com/unedtamps/gobackend/internal/config"
	"github.com/unedtamps/gobackend/pkg/logger"
	"github.com/unedtamps/gobackend/services"
)

func main() {
	log := logger.New("stagging")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	configApp, err := config.NewAppConfiguration(".")
	if err != nil {
		log.Error("Failed Load Configuration", "err", err)
		os.Exit(0)
	}
	log = logger.New(configApp.Server.Environment)

	db, err := database.NewDBInstance(ctx, configApp)
	if err != nil {
		log.Error("Failed Load Database", "err", err)
		log.Info("Close Already Open Database")
		db.Close()
		os.Exit(0)
	}

	server := services.NewServer(db, configApp)
	go func() {
		if err := server.Run(log); err != nil {
			log.Info("Server Error", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Warn("Server Grace Full Shutdown Failed", "error", err)
	}
	db.Close()
	log.Info("Close All Database Connection")
	log.Info("Server is Grace Fully Shutdown")
}
