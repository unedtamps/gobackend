package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/unedtamps/gobackend/internal/bootstrap/database"
	"github.com/unedtamps/gobackend/internal/config"
	"github.com/unedtamps/gobackend/middleware"
)

type Server struct {
	http   *http.Server
	db     *database.DB
	config *config.Config
	engine *gin.Engine
}

type ServerInterface interface {
	Run(log *slog.Logger) error
	Setup(log *slog.Logger)
	MountRoutes() error
	Shutdown(ctx context.Context) error
}

func NewServer(db *database.DB, config *config.Config) ServerInterface {
	if config.Server.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	return &Server{
		http:   &http.Server{},
		db:     db,
		config: config,
		engine: gin.New(),
	}
}

func (s *Server) Run(log *slog.Logger) error {
	s.Setup(log)
	s.MountRoutes()

	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Port),
		Handler: s.engine,
	}

	log.Info(fmt.Sprintf("Server is running in port %d", s.config.Server.Port))
	return s.http.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *Server) MountRoutes() error {
	return nil
}

func (s *Server) Setup(log *slog.Logger) {
	s.engine.Use(gin.Recovery(), middleware.RequestID(), middleware.Logger(log))
}
