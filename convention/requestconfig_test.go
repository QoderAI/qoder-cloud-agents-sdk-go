package convention_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
)

type trackedBody struct {
	io.Reader
	closes int
}

func (b *trackedBody) Close() error { b.closes++; return nil }

type queryParams struct{}

func (queryParams) URLQuery() (url.Values, error) { return url.Values{"q": {"a b"}}, nil }
func requestOptions(fn transportFunc) []option.RequestOption {
	return []option.RequestOption{option.WithBaseURL("https://sdk.test/api/"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: fn})}
}

func TestResponseBodyOwnership(t *testing.T) {
	for _, kind := range []string{"discard", "json", "raw", "error", "invalid_json", "read_error"} {
		t.Run(kind, func(t *testing.T) {
			body := &trackedBody{Reader: strings.NewReader(`{"id":"one"}`)}
			status := 200
			var decoded map[string]any
			var raw *http.Response
			var dst any = &decoded
			switch kind {
			case "discard":
				dst = nil
			case "raw":
				dst = &raw
			case "error":
				status = 503
			case "invalid_json":
				body.Reader = strings.NewReader("{bad")
			case "read_error":
				body.Reader = errorReader{errors.New("read failed")}
			}
			err := convention.ExecuteNewRequest(context.Background(), "GET", "resource", nil, dst, requestOptions(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}}, Body: body, Request: r}, nil
			})...)
			if kind == "error" || kind == "invalid_json" || kind == "read_error" {
				if err == nil {
					t.Fatal("expected failure")
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if kind == "raw" {
				if body.closes != 0 {
					t.Fatal("raw body closed before caller reads")
				}
				b, e := io.ReadAll(raw.Body)
				if e != nil || string(b) != `{"id":"one"}` {
					t.Fatal(string(b), e)
				}
				_ = raw.Body.Close()
			}
			if body.closes != 1 {
				t.Fatalf("response body closed %d times", body.closes)
			}
		})
	}
}

type errorReader struct{ err error }

func (r errorReader) Read([]byte) (int, error) { return 0, r.err }

func TestRequestReplayPreservesReaderPosition(t *testing.T) {
	body := bytes.NewReader([]byte("prefix-payload"))
	_, _ = body.Seek(7, io.SeekStart)
	var sent []string
	cfg, err := convention.NewRequestConfig(context.Background(), "POST", "resource", body, nil, append(requestOptions(func(r *http.Request) (*http.Response, error) {
		b, e := io.ReadAll(r.Body)
		if e != nil {
			t.Fatal(e)
		}
		sent = append(sent, string(b))
		if r.ContentLength != int64(len(b)) {
			t.Errorf("Content-Length=%d, body=%q", r.ContentLength, b)
		}
		status := 429
		if len(sent) > 1 {
			status = 204
		}
		return &http.Response{StatusCode: status, Header: http.Header{"Retry-After-Ms": {"0"}}, Body: io.NopCloser(strings.NewReader("")), Request: r}, nil
	}), option.WithMaxRetries(1))...)
	if err != nil {
		t.Fatal(err)
	}
	if err = cfg.Execute(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sent, []string{"payload", "payload"}) {
		t.Fatalf("replayed bytes changed: %q", sent)
	}
	first, _ := cfg.Request.GetBody()
	second, _ := cfg.Request.GetBody()
	defer first.Close()
	defer second.Close()
	a, _ := io.ReadAll(first)
	b, _ := io.ReadAll(second)
	if string(a) != "payload" || string(b) != "payload" {
		t.Fatalf("replay readers share position: %q %q", a, b)
	}
}
func TestInvalidRequestInputsReturnErrors(t *testing.T) {
	t.Run("invalid_query_url", func(t *testing.T) {
		_, err := convention.NewRequestConfig(context.Background(), "GET", ":bad", queryParams{}, nil)
		if err == nil {
			t.Fatal("accepted invalid URL")
		}
	})
	t.Run("typed_nil_client", func(t *testing.T) {
		var client *http.Client
		_, err := convention.NewRequestConfig(context.Background(), "GET", "resource", nil, nil, option.WithHTTPClient(client))
		if err == nil {
			t.Fatal("accepted nil HTTP client")
		}
	})
	t.Run("nil_response", func(t *testing.T) {
		err := convention.ExecuteNewRequest(context.Background(), "GET", "resource", nil, nil, option.WithBaseURL("https://sdk.test"), option.WithMaxRetries(0), option.WithHTTPClient(nilResponseDoer{}))
		if err == nil {
			t.Fatal("accepted nil response")
		}
	})
}

type nilResponseDoer struct{}

func (nilResponseDoer) Do(*http.Request) (*http.Response, error) { return nil, nil }

func TestRetryCancellationAndBodySafety(t *testing.T) {
	for _, tc := range []struct {
		name, method string
		body         func() io.Reader
		key          string
		status       int
		retryHeader  string
		want         int
	}{
		{"get_503", "GET", nil, "", 503, "", 2},
		{"get_header_opt_out", "GET", nil, "", 503, "false", 1},
		{"post_503", "POST", func() io.Reader { return bytes.NewBufferString("payload") }, "", 503, "true", 1},
		{"post_429", "POST", func() io.Reader { return bytes.NewBufferString("payload") }, "", 429, "", 2},
		{"keyed_post", "POST", func() io.Reader { return bytes.NewBufferString("payload") }, "key", 503, "", 2},
		{"keyed_conflict", "POST", func() io.Reader { return bytes.NewBufferString("payload") }, "key", 409, "true", 1},
		{"unreplayable", "POST", func() io.Reader { return io.LimitReader(strings.NewReader("payload"), 7) }, "key", 429, "", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			var bodies []*trackedBody
			var body io.Reader
			if tc.body != nil {
				body = tc.body()
			}
			opts := append(requestOptions(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Body != nil {
					_, _ = io.ReadAll(r.Body)
					_ = r.Body.Close()
				}
				b := &trackedBody{Reader: strings.NewReader(`{"error":{"message":"retry"}}`)}
				bodies = append(bodies, b)
				h := http.Header{"Retry-After-Ms": {"0"}}
				h.Set("X-Should-Retry", tc.retryHeader)
				if r.Header.Get("X-Qoder-Retry-Count") != string(rune('0'+calls-1)) {
					t.Error("retry count header lost")
				}
				return &http.Response{StatusCode: tc.status, Header: h, Body: b, Request: r}, nil
			}), option.WithMaxRetries(1), option.WithHeader("Idempotency-Key", tc.key))
			if err := convention.ExecuteNewRequest(context.Background(), tc.method, "resource", body, nil, opts...); err == nil {
				t.Fatal("expected final error")
			}
			if calls != tc.want {
				t.Fatalf("calls=%d want=%d", calls, tc.want)
			}
			for _, b := range bodies {
				if b.closes != 1 {
					t.Fatalf("attempt body closes=%d", b.closes)
				}
			}
		})
	}
	t.Run("cancel_during_backoff", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		calls := 0
		opts := append(requestOptions(func(r *http.Request) (*http.Response, error) {
			calls++
			time.AfterFunc(10*time.Millisecond, cancel)
			return &http.Response{StatusCode: 429, Header: http.Header{"Retry-After": {"60"}}, Body: io.NopCloser(strings.NewReader("")), Request: r}, nil
		}), option.WithMaxRetries(2))
		start := time.Now()
		err := convention.ExecuteNewRequest(ctx, "GET", "resource", nil, nil, opts...)
		if !errors.Is(err, context.Canceled) || calls != 1 || time.Since(start) > time.Second {
			t.Fatalf("cancel did not stop backoff: %v, calls=%d", err, calls)
		}
	})
	t.Run("request_timeout", func(t *testing.T) {
		opts := append(requestOptions(func(r *http.Request) (*http.Response, error) { <-r.Context().Done(); return nil, r.Context().Err() }), option.WithRequestTimeout(10*time.Millisecond))
		if err := convention.ExecuteNewRequest(context.Background(), "GET", "resource", nil, nil, opts...); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatal(err)
		}
	})
}
