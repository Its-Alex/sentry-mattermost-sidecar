package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/rpsl/sentry-mattermost-sidecar/config"
	"github.com/tidwall/gjson"

	log "github.com/sirupsen/logrus"
)

// Fields are rendered in Mattermost in this order; empty values are skipped.
type sentryMessage struct {
	Title  string
	Link   string
	Fields [][2]string
}

// The event API url looks like /api/0/projects/<org>/<project>/events/<id>/
var projectSlugFromEventURL = regexp.MustCompile(`/projects/[^/]+/([^/]+)/events/`)

// parseSentryPayload supports two webhook formats:
// - Sentry App (Internal Integration) "event_alert" payload, event is in data.event;
// - legacy WebHooks plugin payload (not available in current self-hosted Sentry), event is in event.
func parseSentryPayload(payload string) sentryMessage {
	if event := gjson.Get(payload, "data.event"); event.Exists() {
		project := ""
		if m := projectSlugFromEventURL.FindStringSubmatch(event.Get("url").String()); m != nil {
			project = m[1]
		}

		server := event.Get("server_name").String()
		if server == "" {
			server = event.Get(`tags.#(0=="server_name").1`).String()
		}

		return sentryMessage{
			Title: event.Get("title").String(),
			Link:  event.Get("web_url").String(),
			Fields: [][2]string{
				{"Culprit", event.Get("culprit").String()},
				{"Project", project},
				{"Environment", event.Get("environment").String()},
				{"Server", server},
				{"Rule", gjson.Get(payload, "data.triggered_rule").String()},
			},
		}
	}

	return sentryMessage{
		Title: gjson.Get(payload, "event.title").String(),
		Link:  gjson.Get(payload, "url").String(),
		Fields: [][2]string{
			{"Culprit", gjson.Get(payload, "culprit").String()},
			{"Project", gjson.Get(payload, "project_slug").String()},
			{"Environment", gjson.Get(payload, "event.environment").String()},
			{"Server", gjson.Get(payload, "event.server_name").String()},
		},
	}
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
		return
	}

	// Create payload for Mattermost
	mmPayload, err := makeMMPayload(channel, sentryJSONPayload)

	if err != nil {
		log.WithError(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"result": "",
			"error":  "internal error",
		})
		return
	}

	// Send webhook request into Mattermost
	cfg := c.Value("cfg").(*config.Config)

	// todo: wrap into goroutine
	resp, err := http.Post(
		cfg.WebhookURL,
		"application/json",
		bytes.NewBuffer(mmPayload),
	)

	if err != nil {
		log.WithError(err).Error("can't send payload to mattermost")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"result": "",
			"error":  "internal error",
		})
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		log.WithFields(log.Fields{
			"channel":    channel,
			"statusCode": resp.StatusCode,
			"response":   string(body),
		}).Error("mattermost rejected payload")
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"result": "",
			"error":  "mattermost rejected payload",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"error":  nil,
	})
}

func makeMMPayload(channel string, sentryJSONPayload string) ([]byte, error) {
	message := parseSentryPayload(sentryJSONPayload)
	fields := makeMMPayloadFields(sentryJSONPayload)

	payload := map[string]interface{}{
		"attachments": []interface{}{
			map[string]interface{}{
				"title":       message.Title,
				"color":       "#FF0000",
				"author_name": "Sentry",
				"author_icon": "https://assets.stickpng.com/images/58482eedcef1014c0b5e4a76.png",
				"title_link":  message.Link,
				"fields":      fields,
			},
		},
	}

	// Without a channel Mattermost posts into the webhook's own channel,
	// which is the only option for webhooks locked to a channel.
	if channel != "" {
		payload["channel"] = channel
	}

	return json.Marshal(payload)
}

func makeMMPayloadFields(sentryJSONPayload string) []interface{} {
	var fields []interface{}

	for _, f := range parseSentryPayload(sentryJSONPayload).Fields {
		if f[1] == "" {
			continue
		}

		fields = append(fields, map[string]interface{}{
			"short": false,
			"title": f[0],
			"value": f[1],
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
