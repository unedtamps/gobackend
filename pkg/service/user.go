package service

import (
	"context"

	"github.com/google/uuid"
)

func (u *user) GetByID(ctx context.Context, Id uuid.UUID) interface{} {
	return nil
}

func (u *user) Haha(ctx context.Context) error {
	return nil
}
