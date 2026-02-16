package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	primary "github.com/unedtamps/gobackend/internal/datastore/primary/gen"
	"github.com/unedtamps/gobackend/pkg/utils"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
	ErrUnauthorized = errors.New("unauthorized access to todo")
)

type Service struct {
	queries primary.Querier
}

func New(queries primary.Querier) *Service {
	return &Service{queries: queries}
}

type Interface interface {
	Create(ctx context.Context, userID utils.ULID, title, description string) (TodoResult, error)
	GetByID(ctx context.Context, id utils.ULID) (TodoResult, error)
	ListByUser(ctx context.Context, userID utils.ULID, limit, offset int32) ([]TodoResult, error)
	Update(ctx context.Context, id, userID utils.ULID, title, description string) (TodoResult, error)
	Delete(ctx context.Context, id, userID utils.ULID) error
	SearchByUser(ctx context.Context, userID utils.ULID, query string, limit, offset int32) ([]TodoResult, error)
}

var _ Interface = (*Service)(nil)

type TodoResult struct {
	ID          utils.ULID
	UserID      utils.ULID
	Title       string
	Description pgtype.Text
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

func mapTodoToResult(todo *primary.Todo) TodoResult {
	return TodoResult{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

func mapTodosToResult(todos []*primary.Todo) []TodoResult {
	results := make([]TodoResult, len(todos))
	for i, todo := range todos {
		results[i] = mapTodoToResult(todo)
	}
	return results
}
