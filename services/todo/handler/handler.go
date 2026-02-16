package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/unedtamps/gobackend/services/todo/service"
)

type Handler struct {
	svc service.Interface
}

func New(svc service.Interface) *Handler {
	return &Handler{svc: svc}
}

type Interface interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Search(c *gin.Context)
}

var _ Interface = (*Handler)(nil)
