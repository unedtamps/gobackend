package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oklog/ulid/v2"
	primary "github.com/unedtamps/gobackend/internal/datastore/primary/gen"
	"github.com/unedtamps/gobackend/pkg/utils"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
)

type Service struct {
	queries primary.Querier
	jwtGen  *utils.JWTGenerator
}

func New(queries primary.Querier, jwtGen *utils.JWTGenerator) *Service {
	return &Service{
		queries: queries,
		jwtGen:  jwtGen,
	}
}

type Interface interface {
	Register(ctx context.Context, email, password string) (UserResult, error)
	Login(ctx context.Context, email, password string) (AuthResult, error)
	GetByID(ctx context.Context, id ulid.ULID) (UserResult, error)
	Update(ctx context.Context, id ulid.ULID, email, password string) (UserResult, error)
	SoftDelete(ctx context.Context, id ulid.ULID) error
}

var _ Interface = (*Service)(nil)

type UserResult struct {
	ID        ulid.ULID
	Email     string
	Status    string
	CreatedAt pgtype.Timestamptz
}

type AuthResult struct {
	User  UserResult
	Token string
}

func mapUserToResult(user *primary.User) UserResult {
	return UserResult{
		ID:        user.ID.ULID,
		Email:     user.Email,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
	}
}
