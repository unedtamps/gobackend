// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package repository

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*GetUserByEmailRow, error)
	MakeTodo(ctx context.Context, arg MakeTodoParams) (uuid.UUID, error)
	QueryTodo(ctx context.Context) ([]*Todolist, error)
	QueryTodoById(ctx context.Context, id uuid.UUID) (*QueryTodoByIdRow, error)
}

var _ Querier = (*Queries)(nil)
