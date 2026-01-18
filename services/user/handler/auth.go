package handler

import (
	"github.com/gin-gonic/gin"
	mysql_gen "github.com/unedtamps/gobackend/internal/datastore/mysql/primary/gen"
	pg_gen "github.com/unedtamps/gobackend/internal/datastore/postgres/primary/gen"
)

func (u *User) LoginUser(c *gin.Context) {
	u.MYSQLPrimary.CreateProduct(c, mysql_gen.CreateProductParams{})

}

func (u *User) RegisterUser(c *gin.Context) {
	u.PGPrimary.CreateUser(c, pg_gen.CreateUserParams{})
}
