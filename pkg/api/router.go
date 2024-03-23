package api

import (
	"errors"
	"net/http"

	mc "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/unedtamps/gobackend/pkg/api/routes"
	"github.com/unedtamps/gobackend/pkg/handler"
	m "github.com/unedtamps/gobackend/pkg/middleware"
	"github.com/unedtamps/gobackend/pkg/repository"
	"github.com/unedtamps/gobackend/pkg/service"
	"github.com/unedtamps/gobackend/util"
)

func (s *Server) SetUpRouter() {
	s.router.Use(mc.Logger)

	s.router.Use(cors.AllowAll().Handler)

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		util.ResponseSuccess(w, nil, 200, "api is working")
	})
	m.SetJwt()

	repo := repository.NewStore(s.db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service, m.TokenAuth)

	s.router.Mount("/user", routes.UserRoutes(handler))
	s.router.Mount("/todo", routes.TodoRoutes(handler))

	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		util.ResponseError(w, 404, errors.New("route not found"))
	})

}
