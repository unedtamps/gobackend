package handler

import (
	"github.com/gin-gonic/gin"
)

func (u *User) LoginUser(c *gin.Context) {
	// u.SECONDARY.CreateProduct(c, secondary.CreateProductParams{})
	c.JSON(200, "success")
}

func (u *User) RegisterUser(c *gin.Context) {
	// u.PRIMARY.CreateUser(c, primary.CreateUserParams{})
	c.JSON(200, "success")
}
