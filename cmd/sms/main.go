package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

		jsonByteData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatalf("Error reading body: %v", err)
		}
		jsonStringData := string(jsonByteData)
		title := gjson.Get(jsonStringData, "event.title").String()

		postBody, err := json.Marshal(map[string]interface{}{
			"channel": channel,
			"attachments": []interface{}{
				map[string]interface{}{
					"title":       title,
					"fallback":    title,
					"color":       "#FF0000",
					"author_name": "Sentry",
					"author_icon": "https://assets.stickpng.com/images/58482eedcef1014c0b5e4a76.png",
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
