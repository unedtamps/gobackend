package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unedtamps/gobackend/middleware"
	"github.com/unedtamps/gobackend/pkg/utils"
	"github.com/unedtamps/gobackend/services/todo/service"
)

func (h *Handler) Create(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	result, err := h.svc.Create(c.Request.Context(), userID, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, TodoResponse{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	todoID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid todo ID",
		})
		return
	}

	result, err := h.svc.GetByID(c.Request.Context(), todoID)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Todo not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	// Verify ownership
	if result.UserID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "Access denied",
		})
		return
	}

	c.JSON(http.StatusOK, TodoResponse{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	})
}

func (h *Handler) List(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req ListTodosRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	results, err := h.svc.ListByUser(c.Request.Context(), userID, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	response := TodoListResponse{
		Data: make([]TodoResponse, len(results)),
	}
	for i, r := range results {
		response.Data[i] = TodoResponse{
			ID:          r.ID,
			Title:       r.Title,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Update(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	todoID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid todo ID",
		})
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	result, err := h.svc.Update(c.Request.Context(), todoID, userID, req.Title, req.Description)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTodoNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Todo not found",
			})
		case errors.Is(err, service.ErrUnauthorized):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Message: "Access denied",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, TodoResponse{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	todoID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid todo ID",
		})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), todoID, userID); err != nil {
		switch {
		case errors.Is(err, service.ErrTodoNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Todo not found",
			})
		case errors.Is(err, service.ErrUnauthorized):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Message: "Access denied",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Message: err.Error(),
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Search(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req SearchTodosRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	results, err := h.svc.SearchByUser(c.Request.Context(), userID, req.Query, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	response := TodoListResponse{
		Data: make([]TodoResponse, len(results)),
	}
	for i, r := range results {
		response.Data[i] = TodoResponse{
			ID:          r.ID,
			Title:       r.Title,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}
