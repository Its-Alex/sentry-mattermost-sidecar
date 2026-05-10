package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/itsalex/sentry-mattermost-sidecar/internal/sentryhook"
)

func init() {
	viper.SetEnvPrefix("sms")

	viper.BindEnv("mattermost_webhook_url")
	viper.BindEnv("host")
	viper.BindEnv("port")

	viper.SetDefault("addr", "0.0.0.0")
	viper.SetDefault("port", "1323")
}

func main() {
	if viper.GetString("mattermost_webhook_url") == "" {
		log.Fatalf("SMS_MATTERMOST_WEBHOOK_URL environment variable must be set!")
	}

	r := sentryhook.NewRouter(sentryhook.RouterConfig{
		WebhookURL: viper.GetString("mattermost_webhook_url"),
	})
	r.Run(fmt.Sprintf(
		"%s:%s",
		viper.GetString("host"),
		viper.GetString("port"),
	))
}
