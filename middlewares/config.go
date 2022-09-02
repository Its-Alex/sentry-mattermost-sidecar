package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rpsl/sentry-mattermost-sidecar/config"
)

func AttachConfig(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo use const instead magick words
		c.Set("cfg", config)
		c.Next()
	}
}
