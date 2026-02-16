package services

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goforj/wire"

	"github.com/unedtamps/gobackend/internal/bootstrap/database"
	"github.com/unedtamps/gobackend/internal/config"
	"github.com/unedtamps/gobackend/middleware"
	"github.com/unedtamps/gobackend/services/user"
)

var ServerSet = wire.NewSet(NewServer, Setup)

type Server struct {
	db     *database.DB
	config *config.Config
	engine *gin.Engine
}

type ServerInterface interface {
	Run(log *slog.Logger) error
}

func MountRoutes(e *gin.Engine, db *database.DB, config *config.Config) {
	userServices := user.InitUserServices(db, config)
	v1 := e.Group("/api/v1")
	{
		user.RegisterRoutes(v1, userServices)
	}
}

func Setup(log *slog.Logger) *gin.Engine {
	e := gin.New()
	e.Use(gin.Recovery(), middleware.RequestID(), middleware.Logger(log))
	return e
}

func NewServer(
	db *database.DB,
	config *config.Config,
	e *gin.Engine,
) ServerInterface {
	if config.Server.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Mount routes here
	MountRoutes(e, db, config)

	return &Server{
		db:     db,
		config: config,
		engine: e,
	}
}

func (s *Server) Run(log *slog.Logger) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Port),
		Handler: s.engine,
	}
	log.Info(fmt.Sprintf("Server is running in port %d", s.config.Server.Port))
	return server.ListenAndServe()
}
