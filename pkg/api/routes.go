package api

import (
	mc "github.com/go-chi/chi/v5/middleware"
	"github.com/unedtamps/gobackend/pkg/handler"
	"github.com/unedtamps/gobackend/pkg/repository"
	"github.com/unedtamps/gobackend/pkg/service"
)

func (s *Server) SetUpRouter() {
	s.router.Use(mc.Logger)

	repo := repository.NewStore(s.db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	s.router.Post("/create", handler.CreateTodo)
	s.router.Get("/todo/{id}", handler.GetTodoByID)
}
