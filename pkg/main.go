package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chi_mid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unedtamps/gobackend/pkg/handler"
	m "github.com/unedtamps/gobackend/pkg/middleware"
	"github.com/unedtamps/gobackend/pkg/repository"
	"github.com/unedtamps/gobackend/pkg/router"
	"github.com/unedtamps/gobackend/pkg/service"
	"github.com/unedtamps/gobackend/utils"
)

type Server struct {
	db     *pgxpool.Pool
	router *chi.Mux
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{
		router: chi.NewRouter(),
		db:     db,
	}
}

func (s *Server) Run() error {
	s.Setup()
	fmt.Println("Server is running on port:", os.Getenv("SERVER_PORT"))
	return http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), s.router)
}

func (s *Server) Setup() {
	s.router.Use(chi_mid.Logger)
	s.router.Use(cors.AllowAll().Handler)

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseSuccess(w, nil, 200, "Golang Backend")
	})
	m.SetJwt()

	repo := repository.NewStore(s.db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	s.router.Mount("/user", router.UserRoutes(handler))
	s.router.Mount("/todo", router.TodoRoutes(handler))

	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseError(w, http.StatusNotFound, errors.New("Route Not Found"))
	})
	s.router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseError(w, http.StatusMethodNotAllowed, errors.New("Method Not Allowed"))
	})

}
