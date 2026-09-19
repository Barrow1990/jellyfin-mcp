package jellyfin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Jellyfin 12 rejects X-MediaBrowser-Token, X-Emby-Token and ?api_key=, so the
// client must authenticate with the Authorization: MediaBrowser scheme.
func TestDoRequestSendsAuthorizationHeader(t *testing.T) {
	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := &JellyfinClient{baseURL: srv.URL, apiKey: "abc123", httpClient: srv.Client()}
	if _, err := c.DoRequest(context.Background(), "GET", "/System/Info", nil, nil); err != nil {
		t.Fatalf("DoRequest: %v", err)
	}

	if want := `MediaBrowser Token="abc123"`; got.Get("Authorization") != want {
		t.Errorf("Authorization = %q, want %q", got.Get("Authorization"), want)
	}
	for _, legacy := range []string{"X-MediaBrowser-Token", "X-Emby-Token"} {
		if v := got.Get(legacy); v != "" {
			t.Errorf("legacy header %s should not be sent, got %q", legacy, v)
		}
	}
}
