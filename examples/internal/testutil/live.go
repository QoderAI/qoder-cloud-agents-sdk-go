//go:build live

package testutil

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
)

func ExecutionTimeout(t *testing.T) time.Duration {
	t.Helper()
	seconds := 180
	if raw := os.Getenv("QODER_E2E_TIMEOUT_SECONDS"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 1800 {
			t.Fatal("QODER_E2E_TIMEOUT_SECONDS must be between 1 and 1800")
		}
		seconds = n
	}
	return time.Duration(seconds) * time.Second
}
func Marker(t *testing.T) string {
	t.Helper()
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(b[:])
}
func RequireE2E(t *testing.T, mode string) {
	t.Helper()
	for _, key := range []string{"PAT", "MODEL", "LIVE_ALLOW_WRITE", "LIVE_ALLOW_EXECUTION"} {
		name := "QODER_" + mode + "_" + key
		v := os.Getenv(name)
		if v == "" || ((key == "LIVE_ALLOW_WRITE" || key == "LIVE_ALLOW_EXECUTION") && v != "true") {
			t.Fatalf("E2E requires %s; no production execution started", name)
		}
	}
}
func RequestLog(t *testing.T) option.RequestOption {
	return option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		start := time.Now()
		res, err := next(r)
		if res != nil {
			t.Logf("http method=%s path=%s status=%d request_id=%s elapsed=%s", r.Method, r.URL.EscapedPath(), res.StatusCode, res.Header.Get("X-Request-Id"), time.Since(start).Round(time.Millisecond))
		} else {
			t.Logf("http method=%s path=%s transport_failed=true", r.Method, r.URL.EscapedPath())
		}
		return res, err
	})
}
func PollPause(ctx context.Context) error {
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return fmt.Errorf("execution deadline: %w", ctx.Err())
	case <-timer.C:
		return nil
	}
}
