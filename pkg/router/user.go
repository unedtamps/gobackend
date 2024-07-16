package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/unedtamps/gobackend/pkg/dto"
	"github.com/unedtamps/gobackend/pkg/handler"
	m "github.com/unedtamps/gobackend/pkg/middleware"
)

func UserRoutes(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(m.Validate[dto.UserRegister])
		r.Post("/register", h.Register)
	})
	r.Group(func(r chi.Router) {
		r.Use(m.Validate[dto.UserLogin])
		r.Post("/login", h.Login)
	})

	return r
}
