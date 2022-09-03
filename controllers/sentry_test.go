package controllers

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func Test_makeMMPayloadFields(t *testing.T) {
	js := getJSONPayloadFromFile()

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

func getJSONPayloadFromFile() string {
	content, err := os.ReadFile("sentry_test_payload.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	return string(content)
}
