package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chi_mid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/unedtamps/gobackend/docs"
	"github.com/unedtamps/gobackend/pkg/api/routes"
	"github.com/unedtamps/gobackend/pkg/handler"
	m "github.com/unedtamps/gobackend/pkg/middleware"
	"github.com/unedtamps/gobackend/pkg/repository"
	"github.com/unedtamps/gobackend/pkg/service"
	"github.com/unedtamps/gobackend/util"
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
	s.SetUpRouter()
	s.SetupDocs()
	fmt.Println("Server is running on port:", os.Getenv("SERVER_PORT"))
	return http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), s.router)
}

func (s *Server) SetUpRouter() {

	s.router.Use(chi_mid.Logger)
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
func (s *Server) SetupDocs() {
	docs.SwaggerInfo.Title = "Golang Boiler Plate Docs"
	docs.SwaggerInfo.Description = "Cards"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
	docs.SwaggerInfo.Schemes = []string{"http", "golang"}
	s.router.Get("/docs/*", httpSwagger.WrapHandler)
}
