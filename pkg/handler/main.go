package handler

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/unedtamps/gobackend/pkg/service"
)

type Handler struct {
	*service.Service
}

func NewHandler(s *service.Service, jwt *jwtauth.JWTAuth) *Handler {
	return &Handler{
		Service: s,
	}
}
