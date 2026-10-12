package convention_test

import (
	"context"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"io"
	"net/http"
	"strings"
	"testing"
)

func authModeClients(t *testing.T, opts ...option.RequestOption) {
	t.Helper()
	f := forward.NewClient(opts...)
	m := managed.NewClient(opts...)
	if _, err := f.Templates.List(context.Background(), forward.TemplateListParams{}); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Agents.List(context.Background(), managed.AgentListParams{}); err != nil {
		t.Fatal(err)
	}
}
func TestAuthenticationModeIsolation(t *testing.T) {
	t.Setenv("QODER_PAT", "ambient-pat")
	t.Setenv("QODER_ACCESS_TOKEN", "")
	want := "Bearer ambient-pat"
	transport := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != want {
			t.Error("incorrect auth", r.Header.Get("Authorization"))
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"data":[]}`)), Request: r}, nil
	})}
	opts := []option.RequestOption{option.WithHTTPClient(transport), option.WithMaxRetries(0)}
	authModeClients(t, opts...)
	want = "Bearer static-sat"
	authModeClients(t, append(opts, option.WithAuthMode("sa"), option.WithAuthToken("static-sat"))...)
	want = "Bearer dynamic"
	authModeClients(t, append(opts, option.WithCredential(&tokenProvider{value: "dynamic"}))...)
	t.Setenv("QODER_ACCESS_TOKEN", "generic-token")
	want = "Bearer generic-token"
	authModeClients(t, append(opts, option.WithAuthMode("sa"))...)
	t.Setenv("QODER_ACCESS_TOKEN", "")
	withoutToken := forward.NewClient(append(opts, option.WithAuthMode("sa"))...)
	_, err := withoutToken.Templates.List(context.Background(), forward.TemplateListParams{})
	if err == nil {
		t.Fatal("SA mode fell back to PAT")
	}
}
func TestSharedSAAuthentication(t *testing.T) {
	t.Setenv("QODER_PAT", "ambient-pat")
	exchanges := 0
	transport := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"data":[]}`
		if strings.HasPrefix(r.URL.Host, "openapi.") {
			exchanges++
			body = `{"access_token":"sat","token_type":"Bearer","expires_in":3600}`
		} else if r.Header.Get("Authorization") != "Bearer sat" {
			t.Fatal("business call did not use SAT")
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})}
	cred, err := convention.NewServiceAccountCredential(convention.ServiceAccountCredentialConfig{SAKey: "key", Region: "global", Scopes: []string{"qca.access", "forward.access"}, HTTPClient: transport})
	if err != nil {
		t.Fatal(err)
	}
	authModeClients(t, option.WithCredential(cred), option.WithHTTPClient(transport))
	if exchanges != 1 {
		t.Fatal("cache not shared", exchanges)
	}
	mixed := managed.NewClient(option.WithCredential(cred), option.WithPAT("personal"))
	_, err = mixed.Agents.List(context.Background(), managed.AgentListParams{})
	if err == nil {
		t.Fatal("mixed auth accepted")
	}
}

func TestAuthenticationHeaderOverridesAndDeletion(t *testing.T) {
	t.Setenv("QODER_PAT", "ambient-pat")
	t.Setenv("QODER_ACCESS_TOKEN", "")
	want := ""
	transport := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != want {
			t.Errorf("authorization = %q, want %q", r.Header.Get("Authorization"), want)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"data":[]}`)), Request: r}, nil
	})}
	base := []option.RequestOption{option.WithHTTPClient(transport), option.WithMaxRetries(0)}
	for _, clear := range []option.RequestOption{option.WithHeaderDel("aUtHoRiZaTiOn"), option.WithHeader("AUTHORIZATION", "")} {
		// Static/environment tokens stay removed; dynamic providers retain legacy injection.
		authModeClients(t, append(base, clear)...)
		authModeClients(t, append(base, option.WithPAT("explicit"), clear)...)
		want = "Bearer dynamic"
		authModeClients(t, append(base, option.WithCredential(&tokenProvider{value: "dynamic"}), clear)...)
		want = ""
	}
	want = "Bearer request-override"
	authModeClients(t, append(base, option.WithCredential(&tokenProvider{err: fmt.Errorf("must not be called")}), option.WithHeader("authorization", want))...)
}

func TestClonedRequestKeepsIndependentAuthentication(t *testing.T) {
	t.Setenv("QODER_PAT", "ambient-pat")
	t.Setenv("QODER_ACCESS_TOKEN", "")
	want := "Bearer dynamic"
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != want {
			t.Errorf("authorization = %q, want %q", r.Header.Get("Authorization"), want)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"data":[]}`)), Request: r}, nil
	})}
	cfg, err := convention.NewRequestConfig(context.Background(), "GET", "https://example.test/", nil, nil,
		convention.EnvironmentAuthentication(), option.WithBaseURL("https://example.test/"), option.WithHTTPClient(httpClient), option.WithMaxRetries(0), option.WithResponseBodyInto(new(any)))
	if err != nil {
		t.Fatal(err)
	}
	copy := cfg.Clone(context.Background())
	if err := copy.Apply(option.WithAuthMode("sa"), option.WithCredential(&tokenProvider{value: "dynamic"})); err != nil {
		t.Fatal(err)
	}
	if err := copy.Execute(); err != nil {
		t.Fatal(err)
	}
	want = "Bearer ambient-pat"
	if err := cfg.Execute(); err != nil {
		t.Fatal(err)
	}
}
