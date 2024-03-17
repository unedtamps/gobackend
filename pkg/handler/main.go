package handler

import "github.com/unedtamps/gobackend/pkg/service"

type Handler struct {
	*service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		Service: s,
	}
}
