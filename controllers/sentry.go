package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rpsl/sentry-mattermost-sidecar/config"
	"github.com/tidwall/gjson"

	log "github.com/sirupsen/logrus"
)

var SentryFields = map[string]string{
	"Culprit":     "culprit",
	"Project":     "project_slug",
	"Environment": "event.environment",
	"Server":      "event.server_name",
}

func SentryWebHookHandler(c *gin.Context) {
	channel := c.Param("channel")

	// Parse payload from Sentry
	sentryJSONPayload, err := getSentryJSONPayload(c.Request.Body)

	if err != nil {
		log.WithError(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"result": "",
			"error":  "can't parse input payload",
		})
	}

	// Create payload for Mattermost
	mmPayload, err := makeMMPayload(channel, sentryJSONPayload)

	if err != nil {
		log.WithError(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"result": "",
			"error":  "internal error",
		})
	}

	// Send webhook request into Mattermost
	cfg := c.Value("cfg").(*config.Config)

	// todo: wrap into goroutine
	resp, err := http.Post(
		cfg.WebhookURL,
		"application/json",
		bytes.NewBuffer(mmPayload),
	)

	defer func() {
		_ = resp.Body.Close()
	}()

	if err != nil {
		log.WithError(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"result": "",
			"error":  "internal error",
		})
	}

	// todo: i think we can bypass response from webhook
	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"error":  nil,
	})
}

func makeMMPayload(channel string, sentryJSONPayload string) ([]byte, error) {
	fields := makeMMPayloadFields(sentryJSONPayload)

	payload := map[string]interface{}{
		"channel": channel,
		"attachments": []interface{}{
			map[string]interface{}{
				"title":       gjson.Get(sentryJSONPayload, "event.title").String(),
				"color":       "#FF0000",
				"author_name": "Sentry",
				"author_icon": "https://assets.stickpng.com/images/58482eedcef1014c0b5e4a76.png",
				"title_link":  gjson.Get(sentryJSONPayload, "url").String(),
				"fields":      fields,
			},
		},
	}

	return json.Marshal(payload)
}

func makeMMPayloadFields(sentryJSONPayload string) []interface{} {
	var fields []interface{}

	for k, v := range SentryFields {
		sVal := gjson.Get(sentryJSONPayload, v).String()

		if sVal == "" {
			continue
		}

		fields = append(fields, map[string]interface{}{
			"short": false,
			"title": k,
			"value": sVal,
		})
	}

	return fields
}

func getSentryJSONPayload(body io.Reader) (string, error) {
	jsonByteData, err := io.ReadAll(body)

	if err != nil {
		log.Errorf("Error reading body: %v", err)
		return "", err
	}

	return string(jsonByteData), err
}
