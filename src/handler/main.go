package handler

import (
	"github.com/unedtamps/gobackend/src/service"
)

type Handler struct {
	*service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		s,
	}
}
