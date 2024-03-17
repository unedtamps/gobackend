package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	return http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), s.router)
}
