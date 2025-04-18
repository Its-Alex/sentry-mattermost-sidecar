package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func init() {
	viper.SetEnvPrefix("sms")

	viper.BindEnv("mattermost_webhook_url")
	viper.BindEnv("host")
	viper.BindEnv("port")

	viper.SetDefault("addr", "0.0.0.0")
	viper.SetDefault("port", "1323")

	if viper.GetString("mattermost_webhook_url") == "" {
		log.Fatalf("SMS_MATTERMOST_WEBHOOK_URL environment variable must be set!")
	}
}

func main() {
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

		var postBody []byte
		if c.Request.Header.Get("Sentry-Hook-Resource") == "error" && gjson.Get(jsonStringData, "action").String() == "created" {
			log.Println("Use error custom integration")

			title := gjson.Get(jsonStringData, "data.error.title").String()

			postBody, err = json.Marshal(map[string]interface{}{
				"channel": channel,
				"attachments": []interface{}{
					map[string]interface{}{
						"title":       title,
						"fallback":    title,
						"color":       "#FF0000",
						"author_name": "Sentry - Errors",
						"author_icon": "https://assets.stickpng.com/images/58482eedcef1014c0b5e4a76.png",
						"author_link": gjson.Get(jsonStringData, "data.error.web_url").String(),
						"title_link":  gjson.Get(jsonStringData, "data.error.web_url").String(),
						"fields": []interface{}{
							map[string]interface{}{
								"short": false,
								"title": "Type",
								"value": gjson.Get(jsonStringData, "action").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Culprit",
								"value": gjson.Get(jsonStringData, "data.error.culprit").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Project ID",
								"value": gjson.Get(jsonStringData, "data.error.project").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Environment",
								"value": gjson.Get(jsonStringData, "data.error.environment").String(),
							},
						},
					},
				},
			})
			if err != nil {
				log.Fatalf("Error during json marshal: %v", err)
			}
		} else if c.Request.Header.Get("Sentry-Hook-Resource") == "event_alert" && gjson.Get(jsonStringData, "action").String() == "triggered" {
			log.Println("Use event_alert custom integration")

			title := gjson.Get(jsonStringData, "data.event_alert.title").String()

			postBody, err = json.Marshal(map[string]interface{}{
				"channel": channel,
				"attachments": []interface{}{
					map[string]interface{}{
						"title":       title,
						"fallback":    title,
						"color":       "#FF0000",
						"author_name": "Sentry - Alert event",
						"author_icon": "https://assets.stickpng.com/images/58482eedcef1014c0b5e4a76.png",
						"author_link": gjson.Get(jsonStringData, "data.event.web_url").String(),
						"title_link":  gjson.Get(jsonStringData, "data.event.web_url").String(),
						"fields": []interface{}{
							map[string]interface{}{
								"short": false,
								"title": "Type",
								"value": gjson.Get(jsonStringData, "action").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Triggered rule",
								"value": gjson.Get(jsonStringData, "data.triggered_rule").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Type",
								"value": gjson.Get(jsonStringData, "data.event.type").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Culprit",
								"value": gjson.Get(jsonStringData, "data.event.culprit").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Project ID",
								"value": gjson.Get(jsonStringData, "data.event.project").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Environment",
								"value": gjson.Get(jsonStringData, "data.event.environment").String(),
							},
						},
					},
				},
			})
			if err != nil {
				log.Fatalf("Error during json marshal: %v", err)
			}
		} else if c.Request.Header.Get("Sentry-Hook-Resource") == "issue" {
			log.Println("Use issue custom integration")

			title := gjson.Get(jsonStringData, "data.issue.title").String()

			if gjson.Get(jsonStringData, "action").String() == "created" {

				postBody, err = json.Marshal(map[string]interface{}{
					"channel": channel,
					"attachments": []interface{}{
						map[string]interface{}{
							"title":       title,
							"fallback":    title,
							"color":       "#FF0000",
							"author_name": "Sentry - Issues",
							"author_icon": "https://assets.stickpng.com/images/58482eedcef1014c0b5e4a76.png",
							// Can't get the link because it's not available in the payload
							// "author_link":
							// "title_link":
							"fields": []interface{}{
								map[string]interface{}{
									"short": false,
									"title": "Type",
									"value": gjson.Get(jsonStringData, "action").String(),
								},
								map[string]interface{}{
									"short": false,
									"title": "Culprit",
									"value": gjson.Get(jsonStringData, "data.issue.culprit").String(),
								},
								map[string]interface{}{
									"short": false,
									"title": "Project ID",
									"value": gjson.Get(jsonStringData, "data.issue.project.id").String(),
								},
								map[string]interface{}{
									"short": false,
									"title": "Project",
									"value": gjson.Get(jsonStringData, "data.issue.project.name").String(),
								},
								map[string]interface{}{
									"short": false,
									"title": "Environment",
									"value": gjson.Get(jsonStringData, "data.issue.environment").String(),
								},
							},
						},
					},
				})
				if err != nil {
					log.Fatalf("Error during json marshal: %v", err)
				}
			}
		} else {
			// Legacy integration
			log.Println("Use legacy integration")

			title := gjson.Get(jsonStringData, "event.title").String()

			postBody, err = json.Marshal(map[string]interface{}{
				"channel": channel,
				"attachments": []interface{}{
					map[string]interface{}{
						"title":       title,
						"fallback":    title,
						"color":       "#FF0000",
						"author_name": "Sentry",
						"author_icon": "https://assets.stickpng.com/images/58482eedcef1014c0b5e4a76.png",
						"author_link": gjson.Get(jsonStringData, "url").String(),
						"title_link":  gjson.Get(jsonStringData, "url").String(),
						"fields": []interface{}{
							map[string]interface{}{
								"short": false,
								"title": "Culprit",
								"value": gjson.Get(jsonStringData, "culprit").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Project",
								"value": gjson.Get(jsonStringData, "project_slug").String(),
							},
							map[string]interface{}{
								"short": false,
								"title": "Environment",
								"value": gjson.Get(jsonStringData, "event.environment").String(),
							},
						},
					},
				},
			})
			if err != nil {
				log.Fatalf("Error during json marshal: %v", err)
			}
		}

		resp, err := http.Post(
			viper.GetString("mattermost_webhook_url"),
			"application/json",
			bytes.NewBuffer(postBody),
		)
		if err != nil {
			log.Fatalf("Error when performing webhook call: %v", err)
		}
		defer resp.Body.Close()
	})

	r.Run(fmt.Sprintf(
		"%s:%s",
		viper.GetString("host"),
		viper.GetString("port"),
	))
}
