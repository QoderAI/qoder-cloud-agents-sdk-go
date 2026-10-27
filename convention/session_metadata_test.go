package convention_test

import (
	"context"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func createMarkedSession(mode string, opts []option.RequestOption, forwardMetadata map[string]any, managedMetadata map[string]string) error {
	if mode == "forward" {
		c := forward.NewClient(opts...)
		_, err := c.Sessions.New(context.Background(), forward.SessionNewParams{IdentityID: "idn", TemplateID: "tmpl", Metadata: forwardMetadata})
		return err
	}
	c := managed.NewClient(opts...)
	_, err := c.Sessions.New(context.Background(), managed.SessionNewParams{Metadata: managedMetadata})
	return err
}

func TestSessionCreationMetadata(t *testing.T) {
	for _, mode := range []string{"forward", "managed"} {
		for _, populated := range []bool{false, true} {
			t.Run(mode+map[bool]string{false: "/omitted", true: "/business"}[populated], func(t *testing.T) {
				var f map[string]any
				var m map[string]string
				if populated {
					f = map[string]any{"ticket": "A-123", "nested": map[string]any{"count": int64(9007199254740993)}}
					m = map[string]string{"ticket": "A-123"}
				}
				beforeF, beforeM := maps.Clone(f), maps.Clone(m)
				calls := 0
				var sentMetadata map[string]string
				opts := append(requestOptions(func(r *http.Request) (*http.Response, error) {
					calls++
					if r.Method != "POST" || !strings.HasSuffix(r.URL.Path, "/sessions") {
						t.Fatalf("unexpected request: %s %s", r.Method, r.URL)
					}
					var body struct {
						Metadata map[string]json.RawMessage `json:"metadata"`
					}
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Fatal(err)
					}
					want := map[string]string{"__qca_sdk_created": "true", "__qca_sdk_language": "go", "__qca_sdk_version": strings.TrimPrefix(r.Header.Get("User-Agent"), "qca-go/")}
					for key, expected := range want {
						var got string
						if err := json.Unmarshal(body.Metadata[key], &got); err != nil || got != expected || got == "" {
							t.Fatalf("%s: %s, want %q", key, body.Metadata[key], expected)
						}
					}
					if populated && string(body.Metadata["ticket"]) != `"A-123"` {
						t.Fatal(body.Metadata)
					}
					if populated && mode == "forward" && string(body.Metadata["nested"]) != `{"count":9007199254740993}` {
						t.Fatal("business metadata changed", body.Metadata)
					}
					sentMetadata = want
					return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"id":"sess_test"}`))}, nil
				}), option.WithPAT("test"), option.WithHeader("X-Qoder-Package-Version", "custom-header"))
				for range 2 {
					if err := createMarkedSession(mode, opts, f, m); err != nil {
						t.Fatal(err)
					}
				}
				if !reflect.DeepEqual(f, beforeF) || !reflect.DeepEqual(m, beforeM) {
					t.Fatal("caller metadata mutated")
				}
				// Reusing exactly matching markers is accepted.
				extra := append(opts, option.WithJSONSet("metadata", sentMetadata))
				if populated {
					extra = append(extra, option.WithJSONSet("metadata.ticket", "A-123"))
					if mode == "forward" {
						extra = append(extra, option.WithJSONSet("metadata.nested", f["nested"]))
					}
				}
				if err := createMarkedSession(mode, extra, nil, nil); err != nil {
					t.Fatal(err)
				}
				if calls != 3 {
					t.Fatal(calls)
				}
			})
		}
	}
}

func TestSessionMetadataConflictsBeforeTransport(t *testing.T) {
	for _, mode := range []string{"forward", "managed"} {
		for _, key := range []string{"__qca_sdk_created", "__qca_sdk_language", "__qca_sdk_version"} {
			for _, value := range []any{"user-value", nil, true, 12} {
				t.Run(mode+"/"+key+"/"+string(mustJSON(t, value)), func(t *testing.T) {
					calls := 0
					opts := append(requestOptions(func(r *http.Request) (*http.Response, error) {
						calls++
						t.Fatal("conflict reached transport")
						return nil, nil
					}), option.WithPAT("test"))
					metadata := map[string]any{key: value}
					for _, override := range []option.RequestOption{
						option.WithJSONSet("metadata", metadata),
						option.WithRequestBody("application/json", mustJSON(t, map[string]any{"metadata": metadata})),
					} {
						err := createMarkedSession(mode, append(opts, override), nil, nil)
						if err == nil || err.Error() != "conflicting SDK metadata field: "+key {
							t.Fatalf("wrong error: %v", err)
						}
					}
					if mode == "forward" {
						if err := createMarkedSession(mode, opts, metadata, nil); err == nil {
							t.Fatal("typed conflict accepted")
						}
					} else if s, ok := value.(string); ok {
						if err := createMarkedSession(mode, opts, nil, map[string]string{key: s}); err == nil {
							t.Fatal("typed conflict accepted")
						}
					}
					if calls != 0 || !reflect.DeepEqual(metadata, map[string]any{key: value}) {
						t.Fatal("unexpected mutation or request")
					}
				})
			}
		}
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestSessionMetadataRawBody(t *testing.T) {
	for _, mode := range []string{"forward", "managed"} {
		for _, raw := range []string{`{}`, `{"metadata":null}`, `{"metadata":{}}`, `{"metadata":{"__qca_sdk_future":"preserved"}}`} {
			calls := 0
			opts := append(requestOptions(func(r *http.Request) (*http.Response, error) {
				calls++
				var body map[string]map[string]string
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body["metadata"]["__qca_sdk_created"] != "true" {
					t.Fatal(body)
				}
				if strings.Contains(raw, "future") && body["metadata"]["__qca_sdk_future"] != "preserved" {
					t.Fatal(body)
				}
				return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
			}), option.WithPAT("test"))
			reader := &trackedBody{Reader: strings.NewReader(raw)}
			if err := createMarkedSession(mode, append(opts, option.WithRequestBody("application/json", reader)), nil, nil); err != nil {
				t.Fatal(err)
			}
			if calls != 1 || reader.closes != 1 {
				t.Fatalf("calls %d, closes %d", calls, reader.closes)
			}
		}
		for _, raw := range []string{`null`, `[]`, `invalid`, `{"metadata":[]}`, `{"metadata":false}`, `{"metadata":"invalid"}`} {
			opts := append(requestOptions(func(r *http.Request) (*http.Response, error) {
				t.Fatal("invalid metadata reached transport")
				return nil, nil
			}), option.WithPAT("test"), option.WithRequestBody("application/json", strings.NewReader(raw)))
			if err := createMarkedSession(mode, opts, nil, nil); err == nil {
				t.Fatal("invalid body accepted:", raw)
			}
		}
	}
}

func TestSessionUpdateDoesNotStampMetadata(t *testing.T) {
	for _, mode := range []string{"forward", "managed"} {
		opts := append(requestOptions(func(r *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(body, map[string]any{"metadata": map[string]any{"__qca_sdk_language": "historic"}}) {
				t.Fatal(body)
			}
			return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
		}), option.WithPAT("test"))
		if mode == "forward" {
			c := forward.NewClient(opts...)
			if _, err := c.Sessions.Update(context.Background(), "sess_test", forward.SessionUpdateParams{Metadata: map[string]any{"__qca_sdk_language": "historic"}}); err != nil {
				t.Fatal(err)
			}
		} else {
			c := managed.NewClient(opts...)
			if _, err := c.Sessions.Update(context.Background(), "sess_test", managed.SessionUpdateParams{Metadata: map[string]any{"__qca_sdk_language": "historic"}}); err != nil {
				t.Fatal(err)
			}
		}
	}
}
