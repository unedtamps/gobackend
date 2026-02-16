package service

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	primary "github.com/unedtamps/gobackend/internal/datastore/primary/gen"
	"github.com/unedtamps/gobackend/pkg/utils"
)

func (s *Service) Create(ctx context.Context, userID utils.ULID, title, description string) (TodoResult, error) {
	var desc pgtype.Text
	if description != "" {
		desc = pgtype.Text{String: description, Valid: true}
	}

	todo, err := s.queries.CreateTodo(ctx, primary.CreateTodoParams{
		ID:          utils.NewULID(),
		UserID:      userID,
		Title:       title,
		Description: desc,
	})
	if err != nil {
		return TodoResult{}, err
	}

	return mapTodoToResult(todo), nil
}

func (s *Service) GetByID(ctx context.Context, id utils.ULID) (TodoResult, error) {
	todo, err := s.queries.GetTodo(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return TodoResult{}, ErrTodoNotFound
		}
		return TodoResult{}, err
	}

	return mapTodoToResult(todo), nil
}

func (s *Service) ListByUser(ctx context.Context, userID utils.ULID, limit, offset int32) ([]TodoResult, error) {
	todos, err := s.queries.GetTodoByUser(ctx, primary.GetTodoByUserParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return mapTodosToResult(todos), nil
}

func (s *Service) Update(ctx context.Context, id, userID utils.ULID, title, description string) (TodoResult, error) {
	// First check if todo exists and belongs to user
	todo, err := s.queries.GetTodo(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return TodoResult{}, ErrTodoNotFound
		}
		return TodoResult{}, err
	}

	if todo.UserID != userID {
		return TodoResult{}, ErrUnauthorized
	}

	var desc pgtype.Text
	if description != "" {
		desc = pgtype.Text{String: description, Valid: true}
	} else {
		desc = todo.Description
	}

	if title == "" {
		title = todo.Title
	}

	updated, err := s.queries.UpdateTodo(ctx, primary.UpdateTodoParams{
		ID:          id,
		Title:       title,
		Description: desc,
	})
	if err != nil {
		return TodoResult{}, err
	}

	return mapTodoToResult(updated), nil
}

func (s *Service) Delete(ctx context.Context, id, userID utils.ULID) error {
	// First check if todo exists and belongs to user
	todo, err := s.queries.GetTodo(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrTodoNotFound
		}
		return err
	}

	if todo.UserID != userID {
		return ErrUnauthorized
	}

	return s.queries.DeleteTodo(ctx, id)
}

func (s *Service) SearchByUser(ctx context.Context, userID utils.ULID, query string, limit, offset int32) ([]TodoResult, error) {
	var searchText pgtype.Text
	if query != "" {
		searchText = pgtype.Text{String: query, Valid: true}
	}

	todos, err := s.queries.SearchTodosByUserAndTitle(ctx, primary.SearchTodosByUserAndTitleParams{
		UserID:  userID,
		Column2: searchText,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}

	return mapTodosToResult(todos), nil
}
