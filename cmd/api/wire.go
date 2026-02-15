//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/goforj/wire"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/unedtamps/gobackend/internal/bootstrap/database"
	"github.com/unedtamps/gobackend/internal/config"
	"github.com/unedtamps/gobackend/services"
)

type App struct {
	Server services.ServerInterface
	Log    *slog.Logger
	DB     *database.DB
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := InitializeApp(ctx, ".")
	if err != nil {
		slog.Default().Error("Error initializing server", "err", err)
		return
	}

	go func() {
		if err := app.Server.Run(app.Log); err != nil {
			app.Log.Info("Server Error", "err", err)
			stop()
		}
	}()
	<-ctx.Done()
	app.Log.Info("Server is shutting down")
	app.DB.Close()
	app.Log.Info("Database connections closed")
	app.Log.Info("Server is gracefully shutdown")
}

func InitializeApp(
	ctx context.Context,
	path string,
) (*App, error) {
	wire.Build(
		config.NewAppConfiguration,
		config.NewLogger,
		database.NewDBInstance,
		services.ServerSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
