package config

import (
	"errors"
	"os"
)

type Config struct {
	WebhookURL string
	ListenPort string
	ListenHost string
	Debug      bool
}

func LoadConfig() (*Config, error) {
	// default values
	config := Config{
		WebhookURL: "",
		ListenPort: "1323",
		ListenHost: "0.0.0.0",
		Debug:      false,
	}

	return readEnvironment(&config)
}

func readEnvironment(config *Config) (*Config, error) {
	if os.Getenv("SMS_MATTERMOST_WEBHOOK_URL") == "" {
		// log.Fatalf("SMS_MATTERMOST_WEBHOOK_URL environment variable must be set!")
		return nil, errors.New("SMS_MATTERMOST_WEBHOOK_URL environment variable must be set")
	}

	config.WebhookURL = os.Getenv("SMS_MATTERMOST_WEBHOOK_URL")

	if os.Getenv("SMS_PORT") != "" {
		config.ListenPort = os.Getenv("SMS_PORT")
	}

	if os.Getenv("SMS_HOST") != "" {
		config.ListenHost = os.Getenv("SMS_HOST")
	}

	if os.Getenv("SMS_DEBUG") != "true" {
		config.Debug = false
	} else {
		config.Debug = true
	}

	return config, nil
}
