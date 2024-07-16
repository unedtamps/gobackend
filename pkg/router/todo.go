package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"github.com/unedtamps/gobackend/pkg/handler"
	m "github.com/unedtamps/gobackend/pkg/middleware"
)

func TodoRoutes(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(m.TokenAuth))
		r.Use(m.Authenticator)
		r.Post("/create", h.CreateTodo)
		r.Get("/todo/{id}", h.GetTodoByID)
	})
	return r
}
