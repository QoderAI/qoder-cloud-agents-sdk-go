package convention

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type saRoundTrip func(*http.Request) (*http.Response, error)

func (f saRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func saResponse(r *http.Request, status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}
}
func saCredential(t *testing.T, handler saRoundTrip) *ServiceAccountCredential {
	t.Helper()
	c, err := NewServiceAccountCredential(ServiceAccountCredentialConfig{SAKey: "secret-key", Region: "global", Scopes: []string{"qca.access", "forward.access"}, HTTPClient: &http.Client{Transport: handler}})
	if err != nil {
		t.Fatal(err)
	}
	return c
}
func TestServiceAccountExchangeCacheAndLifetime(t *testing.T) {
	var calls atomic.Int32
	var seconds atomic.Int64
	c := saCredential(t, func(r *http.Request) (*http.Response, error) {
		n := calls.Add(1)
		if r.Method != "POST" || r.URL.String() != "https://openapi.qoder.sh/api/v1/serviceToken/exchange" || r.Header.Get("Authorization") != "Bearer secret-key" {
			t.Error("incorrect exchange request")
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body["scope"] != "qca.access forward.access" || body["audience"] != "qoder" || body["ttl_seconds"] != float64(3600) {
			t.Error("incorrect exchange body")
		}
		return saResponse(r, 200, fmt.Sprintf(`{"access_token":"sat-%d","token_type":"Bearer","expires_in":"100","expires_at":"2099-01-01T00:00:00Z"}`, n)), nil
	})
	c.now = func() time.Time { return time.Unix(seconds.Load(), 0) }
	token, err := c.Token(context.Background())
	if err != nil || token != "sat-1" {
		t.Fatal(token, err)
	}
	seconds.Store(89)
	token, err = c.Token(context.Background())
	if err != nil || token != "sat-1" || calls.Load() != 1 {
		t.Fatal(token, err, calls.Load())
	}
	seconds.Store(90)
	token, err = c.Token(context.Background())
	if err != nil || token != "sat-2" || calls.Load() != 2 {
		t.Fatal(token, err, calls.Load())
	}
	if strings.Contains(fmt.Sprintf("%#v %+v", c, c), "secret-key") {
		t.Fatal("credential repr leaked key")
	}
}
func TestServiceAccountFailureCooldownAndSafeFallback(t *testing.T) {
	var calls atomic.Int32
	var status atomic.Int32
	status.Store(200)
	var seconds atomic.Int64
	c := saCredential(t, func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		code := int(status.Load())
		body := `{"access_token":"sat","token_type":"Bearer","expires_in":100}`
		if code != 200 {
			body = `{"error":{"code":"unavailable","message":"secret-key"},"request_id":"request-1"}`
		}
		return saResponse(r, code, body), nil
	})
	c.now = func() time.Time { return time.Unix(seconds.Load(), 0) }
	if _, err := c.Token(context.Background()); err != nil {
		t.Fatal(err)
	}
	seconds.Store(90)
	status.Store(503)
	for i := 0; i < 2; i++ {
		token, err := c.Token(context.Background())
		if token != "sat" || err != nil {
			t.Fatal(token, err)
		}
	}
	if calls.Load() != 2 {
		t.Fatal("failure not shared")
	}
	seconds.Store(95)
	_, err := c.Token(context.Background())
	var apiErr *apierror.Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 503 {
		t.Fatal(err)
	}
	seconds.Store(100)
	status.Store(200)
	if _, err := c.Token(context.Background()); err != nil {
		t.Fatal(err)
	}
	seconds.Store(190)
	status.Store(403)
	_, err = c.Token(context.Background())
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 403 || strings.Contains(err.Error(), "secret-key") {
		t.Fatal(err)
	}
	_, err = c.Token(context.Background())
	if err == nil || calls.Load() != 5 {
		t.Fatal("permission failures must not use cached token", calls.Load())
	}
}
func TestServiceAccountConcurrentWaitAndCancellation(t *testing.T) {
	started, release := make(chan struct{}), make(chan struct{})
	var calls atomic.Int32
	c := saCredential(t, func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		close(started)
		<-release
		return saResponse(r, 200, `{"access_token":"sat","token_type":"Bearer","expires_in":3600}`), nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	first := make(chan error, 1)
	go func() { _, err := c.Token(ctx); first <- err }()
	<-started
	cancel()
	if !errors.Is(<-first, context.Canceled) {
		t.Fatal("first waiter did not cancel")
	}
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := c.Token(context.Background())
			if err != nil || token != "sat" {
				t.Error(token, err)
			}
		}()
	}
	close(release)
	wg.Wait()
	if calls.Load() != 1 {
		t.Fatal("concurrent requests exchanged more than once")
	}
}
func TestServiceAccountConcurrentFailureAndRetryAfter(t *testing.T) {
	var calls atomic.Int32
	var seconds atomic.Int64
	c := saCredential(t, func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		res := saResponse(r, 429, `{"error":{"code":"limited"}}`)
		res.Header.Set("Retry-After", "120")
		return res, nil
	})
	c.now = func() time.Time { return time.Unix(seconds.Load(), 0) }
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := c.Token(context.Background())
			if err == nil {
				t.Error("expected rate limit")
			}
		}()
	}
	wg.Wait()
	seconds.Store(119)
	if _, err := c.Token(context.Background()); err == nil {
		t.Fatal("expected rate limit")
	}
	if calls.Load() != 1 {
		t.Fatal("cooldown did not collapse failures")
	}
	seconds.Store(120)
	_, _ = c.Token(context.Background())
	if calls.Load() != 2 {
		t.Fatal("cooldown did not expire")
	}
}
func TestServiceAccountInvalidResponse(t *testing.T) {
	for _, ttl := range []string{`0`, `-1`, `1.5`, `"1.5"`, `true`, `43201`, `null`} {
		t.Run(ttl, func(t *testing.T) {
			c := saCredential(t, func(r *http.Request) (*http.Response, error) {
				return saResponse(r, 200, `{"access_token":"sensitive-token","token_type":"Bearer","expires_in":`+ttl+`}`), nil
			})
			_, err := c.Token(context.Background())
			if err == nil || strings.Contains(err.Error(), "sensitive-token") {
				t.Fatal(err)
			}
		})
	}
}
func TestServiceAccountTimeoutAndRegion(t *testing.T) {
	c := saCredential(t, func(r *http.Request) (*http.Response, error) { <-r.Context().Done(); return nil, r.Context().Err() })
	c.timeout = 10 * time.Millisecond
	_, err := c.Token(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	if c.validateRegion("api.qoder.com.cn") == nil {
		t.Fatal("region mismatch accepted")
	}
	if c.validateRegion("custom.proxy.test") != nil {
		t.Fatal("custom proxy inferred")
	}
}

func TestServiceAccountShortLifetimeAndErrorRedaction(t *testing.T) {
	var seconds atomic.Int64
	c := saCredential(t, func(r *http.Request) (*http.Response, error) {
		seconds.Store(1)
		return saResponse(r, 200, `{"access_token":"sensitive-token","token_type":"Bearer","expires_in":1}`), nil
	})
	c.now = func() time.Time { return time.Unix(seconds.Load(), 0) }
	if _, err := c.Token(context.Background()); err == nil {
		t.Fatal("expired exchange result accepted")
	}
	c = saCredential(t, func(r *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("secret-key sensitive-token")
	})
	if _, err := c.Token(context.Background()); err == nil || strings.Contains(err.Error(), "secret-key") {
		t.Fatal("transport error leaked credentials", err)
	}
	c = saCredential(t, func(r *http.Request) (*http.Response, error) {
		return saResponse(r, 403, `{"error":{"code":"secret-key sensitive-token"},"request_id":"sensitive-token"}`), nil
	})
	c.token, c.safeUntil = "sensitive-token", time.Now().Add(time.Hour)
	_, err := c.Token(context.Background())
	var apiErr *apierror.Error
	if !errors.As(err, &apiErr) || strings.Contains(apiErr.Code+apiErr.RequestID, "sensitive-token") || strings.Contains(apiErr.Code, "secret-key") {
		t.Fatal("cached token leaked in error metadata", err)
	}
}

func TestServiceAccountInvalidResponseFields(t *testing.T) {
	for _, body := range []string{
		`null`, `[]`, `not json`,
		`{"token_type":"Bearer","expires_in":100}`,
		`{"access_token":"token","token_type":"Basic","expires_in":100}`,
		`{"access_token":"token","token_type":"Bearer"}`,
		`{"access_token":"token","token_type":"Bearer","expires_in":100,"expires_at":null}`,
		`{"access_token":"token","token_type":"Bearer","expires_in":100,"expires_at":""}`,
		`{"access_token":"token","token_type":"Bearer","expires_in":100,"expires_at":"2026-10-12"}`,
	} {
		c := saCredential(t, func(r *http.Request) (*http.Response, error) { return saResponse(r, 200, body), nil })
		if _, err := c.Token(context.Background()); err == nil {
			t.Error("invalid exchange result accepted", body)
		}
	}
}
