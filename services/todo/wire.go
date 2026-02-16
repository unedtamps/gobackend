//go:build wireinject
// +build wireinject

package todo

import (
	"github.com/gin-gonic/gin"
	"github.com/goforj/wire"

	"github.com/unedtamps/gobackend/internal/bootstrap/database"
	"github.com/unedtamps/gobackend/internal/config"
	primary "github.com/unedtamps/gobackend/internal/datastore/primary/gen"
	"github.com/unedtamps/gobackend/middleware"
	"github.com/unedtamps/gobackend/pkg/utils"
	"github.com/unedtamps/gobackend/services/todo/handler"
	"github.com/unedtamps/gobackend/services/todo/service"
)

type Features struct {
	Handler handler.Interface
	JWTGen  *utils.JWTGenerator
}

func ProvideQuerier(db *database.DB) primary.Querier {
	return primary.New(db.PRIMARY)
}

func ProvideJWTGen(cfg *config.Config) *utils.JWTGenerator {
	return utils.NewJWTGenerator(cfg.Server.JWTSecret)
}

func ProvideHandler(svc *service.Service) handler.Interface {
	return handler.New(svc)
}

var ProviderSet = wire.NewSet(
	ProvideQuerier,
	ProvideJWTGen,
	service.New,
	ProvideHandler,
	wire.Struct(new(Features), "*"),
)

func InitTodoFeatures(db *database.DB, cfg *config.Config) *Features {
	wire.Build(ProviderSet)
	return nil
}

func RegisterRoutes(r *gin.RouterGroup, s *Features) {
	todos := r.Group("/todos")
	{
		todos.Use(middleware.JWT(s.JWTGen))
		todos.POST("", s.Handler.Create)
		todos.GET("", s.Handler.List)
		todos.GET("/search", s.Handler.Search)
		todos.GET("/:id", s.Handler.GetByID)
		todos.PUT("/:id", s.Handler.Update)
		todos.DELETE("/:id", s.Handler.Delete)
	}
}
