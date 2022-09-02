package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rpsl/sentry-mattermost-sidecar/controllers"
)

func SetRoutes(r *gin.Engine) {
	r.POST("/:channel", controllers.WebHookHandler)

	// todo add NoRoute handler
}
