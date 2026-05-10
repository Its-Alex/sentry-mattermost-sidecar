package sentryhook

import (
	"encoding/json"
	"testing"

	"github.com/tidwall/gjson"
)

func TestField(t *testing.T) {
	fieldMap := field("Type", "created")
	if fieldMap["short"] != false || fieldMap["title"] != "Type" || fieldMap["value"] != "created" {
		t.Fatalf("field: %#v", fieldMap)
	}
	fieldMapEmpty := field("Culprit", "")
	if fieldMapEmpty["value"] != "" {
		t.Fatalf("empty value: %#v", fieldMapEmpty)
	}
	marshaledFieldJSON, err := json.Marshal(field("Label", "é 🔔"))
	if err != nil {
		t.Fatal(err)
	}
	var round map[string]interface{}
	if err := json.Unmarshal(marshaledFieldJSON, &round); err != nil {
		t.Fatal(err)
	}
	if round["short"] != false || round["title"] != "Label" || round["value"] != "é 🔔" {
		t.Fatalf("unicode field after JSON: %#v", round)
	}
}

func TestBuildMattermostJSON_ErrorCreated(t *testing.T) {
	sentryWebhookBody := []byte(`{"action":"created","data":{"error":{"title":"Test error","web_url":"https://sentry.example/issue/1","culprit":"main()","project":"42","environment":"production"}}}`)
	mattermostWebhookJSON, err := BuildMattermostJSON("alerts", "error", sentryWebhookBody)
	if err != nil {
		t.Fatal(err)
	}
	parsedMattermost := gjson.ParseBytes(mattermostWebhookJSON)
	if parsedMattermost.Get("channel").String() != "alerts" {
		t.Fatalf("channel: %s", parsedMattermost.Get("channel").String())
	}
	firstAttachment := parsedMattermost.Get("attachments.0")
	if firstAttachment.Get("title").String() != "Test error" {
		t.Fatalf("title: %s", firstAttachment.Get("title").String())
	}
	if firstAttachment.Get("author_name").String() != "Sentry - Errors" {
		t.Fatalf("author_name: %s", firstAttachment.Get("author_name").String())
	}
	if firstAttachment.Get("author_link").String() != "https://sentry.example/issue/1" {
		t.Fatalf("author_link: %s", firstAttachment.Get("author_link").String())
	}
}

func TestBuildMattermostJSON_EventAlertTriggered(t *testing.T) {
	sentryWebhookBody := []byte(`{"action":"triggered","data":{"event_alert":{"title":"Alert title"},"triggered_rule":"High volume","event":{"web_url":"https://sentry.example/event/1","type":"error","culprit":"handler","project":"7","environment":"staging"}}}`)
	mattermostWebhookJSON, err := BuildMattermostJSON("ops", "event_alert", sentryWebhookBody)
	if err != nil {
		t.Fatal(err)
	}
	parsedMattermost := gjson.ParseBytes(mattermostWebhookJSON)
	if parsedMattermost.Get("channel").String() != "ops" {
		t.Fatalf("channel: %s", parsedMattermost.Get("channel").String())
	}
	firstAttachment := parsedMattermost.Get("attachments.0")
	if firstAttachment.Get("title").String() != "Alert title" {
		t.Fatalf("title: %s", firstAttachment.Get("title").String())
	}
	if firstAttachment.Get("author_name").String() != "Sentry - Alert event" {
		t.Fatalf("author_name: %s", firstAttachment.Get("author_name").String())
	}
	if firstAttachment.Get("fields.1.title").String() != "Triggered rule" || firstAttachment.Get("fields.1.value").String() != "High volume" {
		t.Fatalf("triggered_rule field: %s", firstAttachment.Get("fields").Raw)
	}
}

func TestBuildMattermostJSON_IssueCreated(t *testing.T) {
	sentryWebhookBody := []byte(`{"action":"created","data":{"issue":{"title":"New issue","culprit":"db.Query","environment":"prod","project":{"id":"99","name":"api"}}}}`)
	mattermostWebhookJSON, err := BuildMattermostJSON("bugs", "issue", sentryWebhookBody)
	if err != nil {
		t.Fatal(err)
	}
	firstAttachment := gjson.ParseBytes(mattermostWebhookJSON).Get("attachments.0")
	if firstAttachment.Get("author_name").String() != "Sentry - Issues" {
		t.Fatalf("author_name: %s", firstAttachment.Get("author_name").String())
	}
	if firstAttachment.Get("title").String() != "New issue" {
		t.Fatalf("title: %s", firstAttachment.Get("title").String())
	}
}

func TestBuildMattermostJSON_Legacy(t *testing.T) {
	sentryWebhookBody := []byte(`{"event":{"title":"Legacy event","environment":"dev"},"url":"https://legacy.example/x","culprit":"legacyFn","project_slug":"my-proj"}`)
	mattermostWebhookJSON, err := BuildMattermostJSON("general", "", sentryWebhookBody)
	if err != nil {
		t.Fatal(err)
	}
	parsedMattermost := gjson.ParseBytes(mattermostWebhookJSON)
	if parsedMattermost.Get("channel").String() != "general" {
		t.Fatalf("channel: %s", parsedMattermost.Get("channel").String())
	}
	firstAttachment := parsedMattermost.Get("attachments.0")
	if firstAttachment.Get("author_name").String() != "Sentry" {
		t.Fatalf("author_name: %s", firstAttachment.Get("author_name").String())
	}
	if firstAttachment.Get("title_link").String() != "https://legacy.example/x" {
		t.Fatalf("title_link: %s", firstAttachment.Get("title_link").String())
	}
}

func TestBuildMattermostJSON_IssueNonCreatedNilPayload(t *testing.T) {
	sentryWebhookBody := []byte(`{"action":"resolved","data":{"issue":{"title":"x"}}}`)
	mattermostWebhookJSON, err := BuildMattermostJSON("ch", "issue", sentryWebhookBody)
	if err != nil {
		t.Fatal(err)
	}
	if mattermostWebhookJSON != nil {
		t.Fatalf("expected nil payload, got %q", mattermostWebhookJSON)
	}
}

func TestBuildMattermostJSON_ErrorWrongActionUsesLegacy(t *testing.T) {
	sentryWebhookBody := []byte(`{"action":"resolved","event":{"title":"From legacy shape","environment":"e"},"url":"https://u","culprit":"c","project_slug":"p"}`)
	mattermostWebhookJSON, err := BuildMattermostJSON("ch", "error", sentryWebhookBody)
	if err != nil {
		t.Fatal(err)
	}
	var decodedMattermostPayload map[string]interface{}
	if err := json.Unmarshal(mattermostWebhookJSON, &decodedMattermostPayload); err != nil {
		t.Fatal(err)
	}
	mattermostAttachments, _ := decodedMattermostPayload["attachments"].([]interface{})
	if len(mattermostAttachments) == 0 {
		t.Fatal("expected legacy attachment")
	}
}
