package service

import (
	"context"

	primary "github.com/unedtamps/gobackend/internal/datastore/primary/gen"
	"github.com/unedtamps/gobackend/pkg/utils"
)

func (s *Service) Register(ctx context.Context, email, password string) (UserResult, error) {
	_, err := s.queries.GetUserByEmail(ctx, email)
	if err == nil {
		return UserResult{}, ErrEmailExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return UserResult{}, err
	}
	user, err := s.queries.CreateUser(ctx, primary.CreateUserParams{
		ID:       utils.GenerateULID(),
		Email:    email,
		Password: hashedPassword,
	})
	if err != nil {
		return UserResult{}, err
	}

	return mapUserToResult(user), nil
}

func (s *Service) Login(ctx context.Context, email, password string) (AuthResult, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	if !utils.CheckPassword(password, user.Password) {
		return AuthResult{}, ErrInvalidCredentials
	}

	token, err := s.jwtGen.GenerateToken(utils.Credentials{
		UserID: user.ID.String(),
		Email:  user.Email,
	})
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		User:  mapUserToResult(user),
		Token: token,
	}, nil
}
