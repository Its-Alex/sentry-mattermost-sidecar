package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rpsl/sentry-mattermost-sidecar/controllers"
)

func SetRoutes(r *gin.Engine) {
	// temporary for compatibility
	r.POST("/:channel", controllers.SentryWebHookHandler)

	// todo add NoRoute handler
	// need to use named handlers for potential new features
	r.POST("/sentry/:channel", controllers.SentryWebHookHandler)
}
