package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/unedtamps/gobackend/utils"
)

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	data, err := h.Todo.Testing(context.Background())
	if err != nil {
		utils.ResponseError(w, 500, err)
		return
	}
	utils.ResponseSuccess(w, data, 201, "berhasil memebuat todo")
}

func (h *Handler) GetTodoByID(w http.ResponseWriter, r *http.Request) {

	id := uuid.MustParse(chi.URLParam(r, "id"))

	res := h.Todo.GetByID(context.Background(), id)

	utils.ResponseSuccess(w, res, 200, "berhasi mendapat")
}
