package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unedtamps/gobackend/pkg/repository"
)

type serviceI interface {
	GetByID(ctx context.Context, id uuid.UUID) interface{}
}

type user struct {
	repository.Querier
}

type todo struct {
	repository.Querier
}

type Service struct {
	User user
	Todo todo
}

func NewService(repo repository.Querier) *Service {
	return &Service{
		User: user{repo},
		Todo: todo{repo},
	}
}

type customError struct {
	Code  int
	Error error
}

func newError(code int, err error) *customError {
	return &customError{
		Code:  code,
		Error: err,
	}
}
