package middleware

import (
	"net/http"
	shared "tracking-service/internal"
	service "tracking-service/internal/services"

	"github.com/gin-gonic/gin"
)

func TenantAuthMiddleware(service *service.ApplicationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		apiKey := c.Request.Header.Get("x-api-key")

		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		application, err := service.ValidateAPIKey(ctx, apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(string(shared.TenantApplicationIDKey), application.ID)
		c.Set(string(shared.TenantIDKey), application.TenantID)
		c.Next()
	}
}
