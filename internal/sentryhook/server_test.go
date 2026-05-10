package sentryhook

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func mattermostCaptureServer(t *testing.T) (mattermostWebhookURL string, getLastMattermostPOSTBody func() []byte) {
	t.Helper()
	var lastBodyMutex sync.Mutex
	var lastMattermostPOSTBody []byte
	stubMattermostServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
		}
		lastBodyMutex.Lock()
		lastMattermostPOSTBody = body
		lastBodyMutex.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(stubMattermostServer.Close)
	return stubMattermostServer.URL, func() []byte {
		lastBodyMutex.Lock()
		defer lastBodyMutex.Unlock()
		return bytes.Clone(lastMattermostPOSTBody)
	}
}

func postSentry(t *testing.T, ginEngine *gin.Engine, channel, sentryHookResource, sentryWebhookJSON string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/"+channel, bytes.NewBufferString(sentryWebhookJSON))
	if sentryHookResource != "" {
		req.Header.Set("Sentry-Hook-Resource", sentryHookResource)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ginEngine.ServeHTTP(rec, req)
	return rec
}

// One happy-path HTTP test; payload shapes are covered by TestBuildMattermostJSON_*.
func TestRouter_ForwardsToMattermost(t *testing.T) {
	mattermostWebhookURL, getLastMattermostPOSTBody := mattermostCaptureServer(t)
	ginRouter := NewRouter(RouterConfig{WebhookURL: mattermostWebhookURL})

	sentryWebhookJSON := `{"action":"created","data":{"error":{"title":"Test error","web_url":"https://sentry.example/issue/1","culprit":"main()","project":"42","environment":"production"}}}`
	rec := postSentry(t, ginRouter, "alerts", "error", sentryWebhookJSON)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if gjson.ParseBytes(getLastMattermostPOSTBody()).Get("channel").String() != "alerts" {
		t.Fatalf("forwarded: %s", getLastMattermostPOSTBody())
	}
}

func TestRouter_IssueNonCreatedPostsEmptyBody(t *testing.T) {
	mattermostWebhookURL, getLastMattermostPOSTBody := mattermostCaptureServer(t)
	ginRouter := NewRouter(RouterConfig{WebhookURL: mattermostWebhookURL})

	sentryWebhookJSON := `{"action":"resolved","data":{"issue":{"title":"x"}}}`
	rec := postSentry(t, ginRouter, "ch", "issue", sentryWebhookJSON)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if len(getLastMattermostPOSTBody()) != 0 {
		t.Fatalf("expected empty Mattermost body, got %q", getLastMattermostPOSTBody())
	}
}

func TestRouter_MattermostFailureReturns502(t *testing.T) {
	ginRouter := NewRouter(RouterConfig{
		WebhookURL: "http://127.0.0.1:1",
		HTTPClient: http.DefaultClient,
	})
	rec := postSentry(t, ginRouter, "c", "", `{"event":{"title":"x","environment":"e"},"url":"u","culprit":"c","project_slug":"p"}`)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", rec.Code)
	}
}
