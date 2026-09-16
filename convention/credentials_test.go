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
	t.Setenv("QODER_PAT", "")
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

// Sentinel enforcement is defense-in-depth: without a credential the request must fail
// locally rather than hit the wire with no Authorization header.
func TestMissingCredentialFailsFirstRequest(t *testing.T) {
	t.Setenv("QODER_PAT", "")
	unreachable := &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
		t.Fatalf("transport must not be called without a credential: %s", r.URL.Path)
		return nil, nil
	})}
	m := managed.NewClient(option.WithBaseURL("https://qoder.test/managed/"), option.WithHTTPClient(unreachable), option.WithMaxRetries(0))
	f := forward.NewClient(option.WithBaseURL("https://qoder.test/forward/"), option.WithHTTPClient(unreachable), option.WithMaxRetries(0))
	var noCredentials *convention.NoCredentialsError
	if _, err := m.Agents.List(context.Background(), managed.AgentListParams{}); !errors.As(err, &noCredentials) {
		t.Fatalf("managed: want NoCredentialsError, got %v", err)
	}
	if _, err := f.Templates.List(context.Background(), forward.TemplateListParams{}); !errors.As(err, &noCredentials) {
		t.Fatalf("forward: want NoCredentialsError, got %v", err)
	}
	// Explicit PAT clears the sentinel.
	reachable := &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "Bearer explicit" {
			t.Fatalf("authorization missing: %v", r.Header)
		}
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"data":[],"has_more":false}`)), Request: r}, nil
	})}
	mOK := managed.NewClient(option.WithPAT("explicit"), option.WithBaseURL("https://qoder.test/managed/"), option.WithHTTPClient(reachable), option.WithMaxRetries(0))
	fOK := forward.NewClient(option.WithPAT("explicit"), option.WithBaseURL("https://qoder.test/forward/"), option.WithHTTPClient(reachable), option.WithMaxRetries(0))
	if _, err := mOK.Agents.List(context.Background(), managed.AgentListParams{}); err != nil {
		t.Fatal(err)
	}
	if _, err := fOK.Templates.List(context.Background(), forward.TemplateListParams{}); err != nil {
		t.Fatal(err)
	}
}
