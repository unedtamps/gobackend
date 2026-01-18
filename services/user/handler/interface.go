package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/unedtamps/gobackend/internal/bootstrap/database"
	mysql_gen "github.com/unedtamps/gobackend/internal/datastore/mysql/primary/gen"
	pg_gen "github.com/unedtamps/gobackend/internal/datastore/postgres/primary/gen"
)

type User struct {
	MYSQLPrimary mysql_gen.Querier
	PGPrimary    pg_gen.Querier
}

type UserInterface interface {
	LoginUser(c *gin.Context)
	RegisterUser(c *gin.Context)
}

func NewUserService(db *database.DB) UserInterface {
	return &User{
		MYSQLPrimary: mysql_gen.New(db.PrimaryMySQL()),
		PGPrimary:    pg_gen.New(db.PrimaryPG()),
	}
}
