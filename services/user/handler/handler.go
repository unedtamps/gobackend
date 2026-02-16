package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/unedtamps/gobackend/services/user/service"
)

type Handler struct {
	svc service.Interface
}

func New(svc service.Interface) *Handler {
	return &Handler{svc: svc}
}

type Interface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	DeleteAccount(c *gin.Context)
}

var _ Interface = (*Handler)(nil)
