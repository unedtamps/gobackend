package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/unedtamps/gobackend/services/user/service"
)

func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}
	result, err := h.svc.GetByID(c.Request.Context(), ulid.MustParse(userID.(string)))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        result.ID,
		Email:     result.Email,
		Status:    result.Status,
		CreatedAt: result.CreatedAt,
	})
}

type UpdateProfileRequest struct {
	Email    string `json:"email,omitempty"    binding:"omitempty,email"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	result, err := h.svc.Update(
		c.Request.Context(),
		ulid.MustParse(userID.(string)),
		req.Email,
		req.Password,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        result.ID,
		Email:     result.Email,
		Status:    result.Status,
		CreatedAt: result.CreatedAt,
	})
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	if err := h.svc.SoftDelete(c.Request.Context(), ulid.MustParse(userID.(string))); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
