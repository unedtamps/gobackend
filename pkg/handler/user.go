package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/unedtamps/gobackend/util"
)

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// res := h.Todo.GetByID(context.Background(), uuid.New())
	data, err := h.Todo.Testing(context.Background())
	if err != nil {
		util.ResponseError(w, 500, err)
		return
	}
	util.ResponseSuccess(w, data, "mantep")
}

type Test struct {
	Id string `json:"id"`
}

func (h *Handler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))

	res := h.Todo.GetByID(context.Background(), id)

	util.ResponseSuccess(w, res, "mantep")
}
