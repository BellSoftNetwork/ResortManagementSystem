package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.bellsoft.net/rms/api-core/pkg/auth"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "userID"
	UsernameKey         = "username"
	UserRoleKey         = "userRole"
)

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "인증 헤더가 없습니다")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Unauthorized(c, "잘못된 인증 헤더 형식")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, "토큰이 없습니다")
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "유효하지 않거나 만료된 토큰")
			c.Abort()
			return
		}

		// Parse user ID from subject
		userID, err := strconv.ParseUint(claims.Subject, 10, 32)
		if err != nil {
			response.Unauthorized(c, "토큰의 사용자 ID가 유효하지 않습니다")
			c.Abort()
			return
		}

		c.Set(UserIDKey, uint(userID))
		c.Set(UsernameKey, claims.Username)
		c.Set(UserRoleKey, claims.Authorities)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(UserRoleKey)
		if !exists {
			response.Forbidden(c, "사용자 역할을 찾을 수 없습니다")
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			response.Forbidden(c, "유효하지 않은 사용자 역할")
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "권한이 부족합니다")
		c.Abort()
	}
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(UsernameKey)
	if !exists {
		return "", false
	}

	name, ok := username.(string)
	return name, ok
}

func GetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get(UserRoleKey)
	if !exists {
		return "", false
	}

	role, ok := userRole.(string)
	return role, ok
}
