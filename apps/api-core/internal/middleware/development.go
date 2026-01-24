package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// DevelopmentOnlyMiddleware restricts access to development environments only
func DevelopmentOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		environment := viper.GetString("environment")

		// Only allow in local, development, and staging environments
		if environment != "local" && environment != "development" && environment != "staging" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "This endpoint is not available in production environment",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
