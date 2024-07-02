package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unedtamps/gobackend/pkg/repository"
)

func (t *todo) GetByID(ctx context.Context, id uuid.UUID) interface{} {
	data, _ := t.q.QueryTodoById(ctx, id)
	return data
}

func (t *todo) Testing(ctx context.Context) (interface{}, error) {
	id, err := t.q.MakeTodo(ctx, repository.MakeTodoParams{
		Title:       "mantep",
		Description: "Kerjai sesuatu",
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (t *todo) Create(ctx context.Context) (*uuid.UUID, *customError) {
	return nil, nil
}
