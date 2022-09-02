package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rpsl/sentry-mattermost-sidecar/config"
	"github.com/tidwall/gjson"
)

var SentryFields = map[string]string{
	"Culprit":     "culprit",
	"Project":     "project_slug",
	"Environment": "event.environment",
	"Server":      "event.server_name",
}

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf(err.Error())
	}

	// todo: print motd screen on start
	// log.Println(cfg.WebhookURL)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/:channel", func(c *gin.Context) {
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
	})

	_ = r.Run(fmt.Sprintf(
		"%s:%s",
		cfg.ListenHost,
		cfg.ListenPort,
	))
}
