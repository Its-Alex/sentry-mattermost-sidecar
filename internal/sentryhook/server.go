package sentryhook

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RouterConfig configures the HTTP server that accepts Sentry webhooks.
type RouterConfig struct {
	// WebhookURL is the Mattermost incoming webhook URL.
	WebhookURL string
	// HTTPClient posts to Mattermost; if nil, http.DefaultClient is used.
	HTTPClient *http.Client
}

// NewRouter returns a Gin engine with POST /:channel forwarding to Mattermost.
func NewRouter(cfg RouterConfig) *gin.Engine {
	client := cfg.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/:channel", func(c *gin.Context) {
		channel := c.Param("channel")

		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("read body: %v", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		hookResource := c.Request.Header.Get("Sentry-Hook-Resource")
		postBody, err := BuildMattermostJSON(channel, hookResource, raw)
		if err != nil {
			log.Printf("build mattermost json: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resp, err := client.Post(
			cfg.WebhookURL,
			"application/json",
			bytes.NewBuffer(postBody),
		)
		if err != nil {
			log.Printf("mattermost webhook: %v", err)
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		c.Status(http.StatusOK)
	})

	return r
}
