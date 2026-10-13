package controllers

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

func Test_makeMMPayloadFields(t *testing.T) {
	js := getJSONPayloadFromFile("sentry_test_payload.json")

	type args struct {
		sentryJSONPayload string
	}

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			"Positive",
			args{
				sentryJSONPayload: js,
			},
			[]interface{}{
				map[string]interface{}{
					"short": false,
					"title": "Culprit",
					"value": "/vendor/laravel/framework/src/Illuminate/Container/BoundMethod.php in Illuminate\\Container\\BoundMethod::Illuminate\\Container\\{closure}",
				},
				map[string]interface{}{
					"title": "Project",
					"value": "platform",
					"short": false,
				},
				map[string]interface{}{
					"short": false,
					"title": "Environment",
					"value": "production",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeMMPayloadFields(tt.args.sentryJSONPayload); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeMMPayloadFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getJSONPayloadFromFile(name string) string {
	content, err := os.ReadFile(name)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	return string(content)
}

func Test_makeMMPayloadFields_SentryApp(t *testing.T) {
	js := getJSONPayloadFromFile("sentry_app_test_payload.json")

	want := []interface{}{
		map[string]interface{}{
			"short": false,
			"title": "Culprit",
			"value": "App\\Http\\Controllers\\ParkingController@show",
		},
		map[string]interface{}{
			"short": false,
			"title": "Project",
			"value": "backend",
		},
		map[string]interface{}{
			"short": false,
			"title": "Environment",
			"value": "production",
		},
		map[string]interface{}{
			"short": false,
			"title": "Server",
			"value": "parkly-prod-01",
		},
		map[string]interface{}{
			"short": false,
			"title": "Rule",
			"value": "New issue -> mattermost",
		},
	}

	if got := makeMMPayloadFields(js); !reflect.DeepEqual(got, want) {
		t.Errorf("makeMMPayloadFields() = %v, want %v", got, want)
	}

	message := parseSentryPayload(js)

	if message.Title != "ErrorException: Undefined array key \"id\"" {
		t.Errorf("Title = %q", message.Title)
	}

	if message.Link != "https://sentry.example.team/organizations/example/issues/10/events/7ad3823827c244b499c17c914e2d5736/" {
		t.Errorf("Link = %q", message.Link)
	}
}

func Test_makeMMPayload_Channel(t *testing.T) {
	js := getJSONPayloadFromFile("sentry_app_test_payload.json")

	tests := []struct {
		name    string
		channel string
		want    interface{}
	}{
		{"Webhook channel", "", nil},
		{"Explicit channel", "backend", "backend"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := makeMMPayload(tt.channel, js)
			if err != nil {
				t.Fatal(err)
			}

			var got map[string]interface{}
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatal(err)
			}

			if got["channel"] != tt.want {
				t.Errorf("channel = %v, want %v", got["channel"], tt.want)
			}
		})
	}
}
