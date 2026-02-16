package service

import (
	"context"

	primary "github.com/unedtamps/gobackend/internal/datastore/primary/gen"
	"github.com/unedtamps/gobackend/pkg/utils"
)

func (s *Service) GetByID(ctx context.Context, id utils.ULID) (UserResult, error) {
	user, err := s.queries.GetUser(ctx, id)
	if err != nil {
		return UserResult{}, ErrUserNotFound
	}
	return mapUserToResult(user), nil
}

func (s *Service) Update(
	ctx context.Context,
	id utils.ULID,
	email, password string,
) (UserResult, error) {
	var hashedPassword string
	if password != "" {
		var err error
		hashedPassword, err = utils.HashPassword(password)
		if err != nil {
			return UserResult{}, err
		}
	}

	user, err := s.queries.UpdateUser(ctx, primary.UpdateUserParams{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
	})
	if err != nil {
		return UserResult{}, err
	}

	return mapUserToResult(user), nil
}

func (s *Service) SoftDelete(ctx context.Context, id utils.ULID) error {
	_, err := s.queries.SoftDeleteUser(ctx, id)
	return err
}
