package handler

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/unedtamps/gobackend/pkg/utils"
)

type RegisterRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Email    string `json:"email,omitempty"    binding:"omitempty,email"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}

type UserResponse struct {
	ID        utils.ULID         `json:"id"`
	Email     string             `json:"email"`
	Status    string             `json:"status"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
