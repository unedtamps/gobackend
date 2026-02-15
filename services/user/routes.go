package user

import (
	"github.com/gin-gonic/gin"
	"github.com/unedtamps/gobackend/internal/bootstrap/database"
	primary "github.com/unedtamps/gobackend/internal/datastore/primary/gen"
	secondary "github.com/unedtamps/gobackend/internal/datastore/secondary/gen"
	"github.com/unedtamps/gobackend/services/user/handler"
)

func NewUserService(db *database.DB) handler.UserInterface {
	return &handler.User{
		PRIMARY:   primary.New(db.PRIMARY),
		SECONDARY: secondary.New(db.SECONDARY),
	}
}

func UserRoutes(r *gin.RouterGroup, handler handler.UserInterface) {
	r.GET("/register", handler.LoginUser)
	r.GET("/login", handler.RegisterUser)
}
