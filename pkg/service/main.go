package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unedtamps/gobackend/pkg/repository"
)

type serviceI interface {
	GetByID(ctx context.Context, id uuid.UUID) interface{}
}

type Service struct {
	User userServiceI
	Todo todoServiceI
}

func NewService(repo repository.Querier) *Service {
	return &Service{
		User: &user{repo},
		Todo: &todo{repo},
	}
}

// user
type userServiceI interface {
	serviceI
	Haha(ctx context.Context) error
}

type user struct {
	repository.Querier
}

// todo
type todoServiceI interface {
	serviceI
	Testing(ctx context.Context) (interface{}, error)
}
type todo struct {
	repository.Querier
}
