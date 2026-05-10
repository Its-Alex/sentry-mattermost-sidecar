package sentryhook

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

const sentryAuthorIcon = "https://assets.stickpng.com/images/58482eedcef1014c0b5e4a76.png"

// BuildMattermostJSON transforms a Sentry webhook body into Mattermost incoming-webhook JSON.
// hookResource is the Sentry-Hook-Resource header value (may be empty for legacy payloads).
// It returns (nil, nil) when no Mattermost message should be sent (e.g. issue hook without action "created").
func BuildMattermostJSON(channel, hookResource string, sentryJSON []byte) ([]byte, error) {
	s := string(sentryJSON)
	action := gjson.Get(s, "action").String()

	switch {
	case hookResource == "error" && action == "created":
		title := gjson.Get(s, "data.error.title").String()
		return json.Marshal(map[string]interface{}{
			"channel": channel,
			"attachments": []interface{}{
				map[string]interface{}{
					"title":       title,
					"fallback":    title,
					"color":       "#FF0000",
					"author_name": "Sentry - Errors",
					"author_icon": sentryAuthorIcon,
					"author_link": gjson.Get(s, "data.error.web_url").String(),
					"title_link":  gjson.Get(s, "data.error.web_url").String(),
					"fields": []interface{}{
						field("Type", gjson.Get(s, "action").String()),
						field("Culprit", gjson.Get(s, "data.error.culprit").String()),
						field("Project ID", gjson.Get(s, "data.error.project").String()),
						field("Environment", gjson.Get(s, "data.error.environment").String()),
					},
				},
			},
		})

	case hookResource == "event_alert" && action == "triggered":
		title := gjson.Get(s, "data.event_alert.title").String()
		return json.Marshal(map[string]interface{}{
			"channel": channel,
			"attachments": []interface{}{
				map[string]interface{}{
					"title":       title,
					"fallback":    title,
					"color":       "#FF0000",
					"author_name": "Sentry - Alert event",
					"author_icon": sentryAuthorIcon,
					"author_link": gjson.Get(s, "data.event.web_url").String(),
					"title_link":  gjson.Get(s, "data.event.web_url").String(),
					"fields": []interface{}{
						field("Type", gjson.Get(s, "action").String()),
						field("Triggered rule", gjson.Get(s, "data.triggered_rule").String()),
						field("Type", gjson.Get(s, "data.event.type").String()),
						field("Culprit", gjson.Get(s, "data.event.culprit").String()),
						field("Project ID", gjson.Get(s, "data.event.project").String()),
						field("Environment", gjson.Get(s, "data.event.environment").String()),
					},
				},
			},
		})

	case hookResource == "issue":
		if action != "created" {
			return nil, nil
		}
		title := gjson.Get(s, "data.issue.title").String()
		return json.Marshal(map[string]interface{}{
			"channel": channel,
			"attachments": []interface{}{
				map[string]interface{}{
					"title":       title,
					"fallback":    title,
					"color":       "#FF0000",
					"author_name": "Sentry - Issues",
					"author_icon": sentryAuthorIcon,
					"fields": []interface{}{
						field("Type", gjson.Get(s, "action").String()),
						field("Culprit", gjson.Get(s, "data.issue.culprit").String()),
						field("Project ID", gjson.Get(s, "data.issue.project.id").String()),
						field("Project", gjson.Get(s, "data.issue.project.name").String()),
						field("Environment", gjson.Get(s, "data.issue.environment").String()),
					},
				},
			},
		})

	default:
		title := gjson.Get(s, "event.title").String()
		return json.Marshal(map[string]interface{}{
			"channel": channel,
			"attachments": []interface{}{
				map[string]interface{}{
					"title":       title,
					"fallback":    title,
					"color":       "#FF0000",
					"author_name": "Sentry",
					"author_icon": sentryAuthorIcon,
					"author_link": gjson.Get(s, "url").String(),
					"title_link":  gjson.Get(s, "url").String(),
					"fields": []interface{}{
						field("Culprit", gjson.Get(s, "culprit").String()),
						field("Project", gjson.Get(s, "project_slug").String()),
						field("Environment", gjson.Get(s, "event.environment").String()),
					},
				},
			},
		})
	}
}

func field(title, value string) map[string]interface{} {
	return map[string]interface{}{
		"short": false,
		"title": title,
		"value": value,
	}
}
