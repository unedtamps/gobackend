package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unedtamps/gobackend/src/dto"
	"github.com/unedtamps/gobackend/src/repository"
	"github.com/unedtamps/gobackend/utils"
)

func (u *user) GetByID(ctx context.Context, Id uuid.UUID) interface{} {
	return nil
}

func (u *user) Haha(ctx context.Context) error {
	return nil
}

func (u *user) RegisterUser(
	ctx context.Context,
	data dto.UserRegister,
) (*uuid.UUID, *customError) {
	hashed, err := utils.GenereateHash(data.Password)
	if err != nil {
		return nil, newError(500, err)
	}
	id, err := u.q.CreateUser(ctx, repository.CreateUserParams{
		Email:    data.Email,
		Password: hashed,
	})
	if err != nil {
		return nil, newError(500, err)
	}
	return &id, nil
}

func (u *user) LoginUser(ctx context.Context, data dto.UserLogin) (*uuid.UUID, *customError) {
	res, err := u.q.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, newError(404, err)
	}

	err = utils.CompareHash(data.Password, res.Password)
	if err != nil {
		return nil, newError(401, err)
	}
	return &res.ID, nil
}
