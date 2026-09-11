package convention_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

type tokenProvider struct {
	value string
	err   error
}

func (p *tokenProvider) Token(context.Context) (string, error) { return p.value, p.err }

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestCredentialSharedByBothModes(t *testing.T) {
	t.Setenv("QODER_ACCESS_TOKEN", "")
	credential := &tokenProvider{value: "first"}
	var shared convention.Credential = credential
	calls := 0
	client := &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Header.Get("Authorization") != "Bearer "+credential.value {
			t.Fatal("credential not applied", r.Header)
		}
		if r.URL.Path != "/managed/agents" && r.URL.Path != "/forward/templates" {
			t.Fatal(r.URL.Path)
		}
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"data":[],"has_more":false}`)), Request: r}, nil
	})}
	m := managed.NewClient(option.WithCredential(shared), option.WithBaseURL("https://qoder.test/managed/"), option.WithHTTPClient(client), option.WithMaxRetries(0))
	f := forward.NewClient(option.WithCredential(shared), option.WithBaseURL("https://qoder.test/forward"), option.WithHTTPClient(client), option.WithMaxRetries(0))
	for _, token := range []string{"first", "rotated"} {
		credential.value = token
		if _, err := m.Agents.List(context.Background(), managed.AgentListParams{}); err != nil {
			t.Fatal(err)
		}
		if _, err := f.Templates.List(context.Background(), forward.TemplateListParams{}); err != nil {
			t.Fatal(err)
		}
	}
	credential.err = errors.New("token unavailable")
	if _, err := m.Agents.List(context.Background(), managed.AgentListParams{}); !errors.Is(err, credential.err) {
		t.Fatal(err)
	}
	if _, err := f.Templates.List(context.Background(), forward.TemplateListParams{}); !errors.Is(err, credential.err) {
		t.Fatal(err)
	}
	if calls != 4 {
		t.Fatalf("token errors must not reach transport: %d calls", calls)
	}
}
