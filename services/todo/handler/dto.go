package handler

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/unedtamps/gobackend/pkg/utils"
)

type CreateTodoRequest struct {
	Title       string `json:"title"       binding:"required,max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title,omitempty"       binding:"omitempty,max=255"`
	Description string `json:"description,omitempty" binding:"omitempty,max=1000"`
}

type ListTodosRequest struct {
	Limit  int32 `form:"limit,default=10" binding:"min=1,max=100"`
	Offset int32 `form:"offset,default=0" binding:"min=0"`
}

type SearchTodosRequest struct {
	Query  string `form:"q"                binding:"max=255"`
	Limit  int32  `form:"limit,default=10" binding:"min=1,max=100"`
	Offset int32  `form:"offset,default=0" binding:"min=0"`
}

type TodoResponse struct {
	ID          utils.ULID         `json:"id"`
	Title       string             `json:"title"`
	Description pgtype.Text        `json:"description"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type TodoListResponse struct {
	Data  []TodoResponse `json:"data"`
	Total int64          `json:"total,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
