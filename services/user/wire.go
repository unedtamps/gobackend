//go:build wireinject
// +build wireinject

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/goforj/wire"
	"github.com/unedtamps/gobackend/internal/bootstrap/database"
	"github.com/unedtamps/gobackend/internal/config"
	primary "github.com/unedtamps/gobackend/internal/datastore/primary/gen"
	"github.com/unedtamps/gobackend/middleware"
	"github.com/unedtamps/gobackend/pkg/utils"
	"github.com/unedtamps/gobackend/services/user/handler"
	"github.com/unedtamps/gobackend/services/user/service"
)

type Dependencies struct {
	Queries   primary.Querier
	JWTSecret string
}

type Services struct {
	h      handler.Interface
	jwtGen *utils.JWTGenerator
}

func newDependencies(db *database.DB, config *config.Config) Dependencies {
	queries := primary.New(db.PRIMARY)
	return Dependencies{
		Queries:   queries,
		JWTSecret: config.Server.JWTSecret,
	}
}

func newRoutes(deps Dependencies) *Services {
	jwtGen := utils.NewJWTGenerator(deps.JWTSecret)
	svc := service.New(deps.Queries, jwtGen)
	return &Services{
		h:      handler.New(svc),
		jwtGen: jwtGen,
	}
}

func RegisterRoutes(r *gin.RouterGroup, u *Services) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", u.h.Register)
		auth.POST("/login", u.h.Login)
	}

	user := r.Group("/user")
	{
		user.Use(middleware.JWT(u.jwtGen))
		user.GET("/profile", u.h.GetProfile)
		user.PUT("/profile", u.h.UpdateProfile)
		user.DELETE("/account", u.h.DeleteAccount)
	}
}

func InitUserServices(
	db *database.DB,
	config *config.Config,
) *Services {
	wire.Build(
		newDependencies,
		newRoutes,
	)
	return &Services{}
}
