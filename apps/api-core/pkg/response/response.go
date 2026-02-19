package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SingleResponse[T any] struct {
	Value T `json:"value"`
}

type ListResponse[T any] struct {
	Values []T                    `json:"values"`
	Page   map[string]interface{} `json:"page"`
	Filter interface{}            `json:"filter"`
}

// Spring Boot 호환 에러 응답
type ErrorResponse struct {
	Message     string   `json:"message"`
	Errors      []string `json:"errors"`
	FieldErrors []string `json:"fieldErrors"`
}

type Meta struct {
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"requestId,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Version    string      `json:"version,omitempty"`
}

type Pagination struct {
	Page          int   `json:"page"`
	Size          int   `json:"size"`
	TotalPages    int   `json:"totalPages"`
	TotalElements int64 `json:"totalElements"`
}

func Success[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, SingleResponse[T]{
		Value: data,
	})
}

func SuccessList[T any](c *gin.Context, data []T, pagination *Pagination) {
	page := map[string]interface{}{
		"index":         pagination.Page, // 이미 0-based로 받음
		"size":          pagination.Size,
		"totalElements": pagination.TotalElements,
		"totalPages":    pagination.TotalPages,
	}

	c.JSON(http.StatusOK, ListResponse[T]{
		Values: data,
		Page:   page,
		Filter: map[string]interface{}{},
	})
}

func SuccessListWithFilter[T any](c *gin.Context, data []T, pagination *Pagination, filter interface{}) {
	page := map[string]interface{}{
		"index":         pagination.Page, // 이미 0-based로 받음
		"size":          pagination.Size,
		"totalElements": pagination.TotalElements,
		"totalPages":    pagination.TotalPages,
	}

	c.JSON(http.StatusOK, ListResponse[T]{
		Values: data,
		Page:   page,
		Filter: filter,
	})
}

func Created[T any](c *gin.Context, data T) {
	c.JSON(http.StatusCreated, SingleResponse[T]{
		Value: data,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, statusCode int, message string, errors []string, fieldErrors []string) {
	c.JSON(statusCode, ErrorResponse{
		Message:     message,
		Errors:      errors,
		FieldErrors: fieldErrors,
	})
}

func BadRequest(c *gin.Context, message string, details ...string) {
	Error(c, http.StatusBadRequest, message, nil, details)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil, nil)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil, nil)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil, nil)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message, nil, nil)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message, nil, nil)
}

func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, message, nil, nil)
}
