package live

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestExampleAuthenticationSelection(t *testing.T) {
	for _, mode := range []string{"pat", "sa"} {
		t.Run(mode, func(t *testing.T) {
			values := map[string]string{"QODER_AUTH_MODE": mode, "QODER_PAT": "ambient-pat", "QODER_ACCESS_TOKEN": "static-sat"}
			config := Config{Mode: "forward"}
			if err := config.configureAuthentication(func(key string) string { return values[key] }); err != nil {
				t.Fatal(err)
			}
			want := "Bearer ambient-pat"
			if mode == "sa" {
				want = "Bearer static-sat"
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("Authorization"); got != want {
					t.Errorf("Authorization = %q, want %q", got, want)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[]}`))
			}))
			defer server.Close()
			config.BaseURL = server.URL
			client := forward.NewClient(append(config.Options(), option.WithHTTPClient(server.Client()))...)
			if _, err := client.Models.List(context.Background()); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSAExampleAuthenticationDoesNotFallBackToPAT(t *testing.T) {
	for _, values := range []map[string]string{
		{"QODER_AUTH_MODE": "sa", "QODER_PAT": "ambient-pat"},
		{"QODER_AUTH_MODE": "sa", "QODER_SA_KEY": "key", "QODER_ACCESS_TOKEN": "sat"},
		{"QODER_AUTH_MODE": "invalid", "QODER_PAT": "ambient-pat"},
	} {
		config := Config{Mode: "forward"}
		if err := config.configureAuthentication(func(key string) string { return values[key] }); err == nil {
			t.Fatal("invalid example authentication unexpectedly accepted")
		}
	}
}
