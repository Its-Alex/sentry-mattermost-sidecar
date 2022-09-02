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

func WebHookHandler(c *gin.Context) {
	channel := c.Param("channel")

	jsonByteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}
	jsonStringData := string(jsonByteData)

	var fields []interface{}

	for k, v := range SentryFields {
		sVal := gjson.Get(jsonStringData, v).String()

		if sVal == "" {
			continue
		}

		fields = append(fields, map[string]interface{}{
			"short": false,
			"title": k,
			"value": sVal,
		})
	}

	payload := map[string]interface{}{
		"channel": channel,
		"attachments": []interface{}{
			map[string]interface{}{
				"title":       gjson.Get(jsonStringData, "event.title").String(),
				"color":       "#FF0000",
				"author_name": "Sentry",
				"author_icon": "https://assets.stickpng.com/images/58482eedcef1014c0b5e4a76.png",
				"title_link":  gjson.Get(jsonStringData, "url").String(),
				"fields":      fields,
			},
		},
	}

	mmPayload, err := json.Marshal(payload)

	if err != nil {
		log.Fatalf("Error during json marshal: %v", err)
	}

	cfg := c.Value("cfg").(*config.Config)

	resp, err := http.Post(
		cfg.WebhookURL,
		"application/json",
		bytes.NewBuffer(mmPayload),
	)

	if err != nil {
		log.Fatalf("Error when performing webhook call: %v", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()
}
