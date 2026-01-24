package middleware

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
	"gorm.io/gorm"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			handleError(c, err.Err)
		}
	}
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.NotFound(c, "존재하지 않는 데이터")
		return
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		response.Conflict(c, "이미 존재하는 데이터")
		return
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fieldErrors := make([]string, 0, len(validationErrors))
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			// Spring Boot 스타일의 에러 메시지
			msg := fmt.Sprintf("'%s'은(는) %s (요청 값: %v)", field, getValidationMessage(tag, e.Param()), e.Value())
			fieldErrors = append(fieldErrors, msg)
		}
		response.BadRequest(c, "잘못된 요청", fieldErrors...)
		return
	}

	if c.Writer.Status() >= 400 {
		return
	}

	response.InternalServerError(c, "서버 오류")
}

func getValidationMessage(tag string, param string) string {
	switch tag {
	case "required":
		return "필수 값입니다"
	case "min":
		return fmt.Sprintf("최소 %s 이상이어야 합니다", param)
	case "max":
		return fmt.Sprintf("최대 %s 이하여야 합니다", param)
	case "email":
		return "유효한 이메일 형식이 아닙니다"
	case "oneof":
		return fmt.Sprintf("다음 중 하나여야 합니다: %s", param)
	default:
		return fmt.Sprintf("%s 검증에 실패했습니다", tag)
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			response.InternalServerError(c, err)
		} else {
			response.InternalServerError(c, "서버 오류")
		}
		c.Abort()
	})
}
