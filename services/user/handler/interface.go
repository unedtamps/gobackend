package handler

import (
	"github.com/gin-gonic/gin"
	primary "github.com/unedtamps/gobackend/internal/datastore/primary/gen"
	secondary "github.com/unedtamps/gobackend/internal/datastore/secondary/gen"
)

type User struct {
	PRIMARY   primary.Querier
	SECONDARY secondary.Querier
}

type UserInterface interface {
	LoginUser(c *gin.Context)
	RegisterUser(c *gin.Context)
}
