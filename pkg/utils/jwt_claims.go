package utils

import (
	"context"
	"encoding/json"

	"github.com/go-chi/jwtauth/v5"
)

func JwtClaims[T any](ctx context.Context) (*T, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	var res T
	claimsjson, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(claimsjson), &res); err != nil {
		return nil, err
	}
	return &res, nil
}
