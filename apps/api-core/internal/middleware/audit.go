package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.bellsoft.net/rms/api-core/internal/audit"
)

// AuditMiddleware sets user context for audit logging
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user information from the context (set by AuthMiddleware)
		var userID *uint
		var username string

		if id, exists := GetUserID(c); exists {
			userID = &id
		}

		if name, exists := GetUsername(c); exists {
			username = name
		}

		// Set audit context with user information
		ctx := audit.SetUserContext(c.Request.Context(), userID, username)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
