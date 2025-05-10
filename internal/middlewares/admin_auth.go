package middleware

import (
	"net/http"

	shared "tracking-service/internal"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AdminAuthMiddleware(config *shared.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		apiKey := c.Request.Header.Get("x-api-key")

		if apiKey == "" {
			log.WithContext(ctx).Warnf("Invalid API key '%s' for path %v", apiKey, c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if apiKey != config.AdminApiKey {
			log.WithContext(ctx).Warnf("Invalid API key '%s' for path %v", apiKey, c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}
