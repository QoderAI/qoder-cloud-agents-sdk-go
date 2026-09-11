//go:build live

package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	shared "github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/testutil"
)

func ExecutionTimeout(t *testing.T) time.Duration {
	t.Helper()
	return shared.ExecutionTimeout(t)
}

func Marker(t *testing.T) string {
	t.Helper()
	return shared.Marker(t)
}

func RequireE2E(t *testing.T, mode string) {
	t.Helper()
	shared.RequireE2E(t, mode)
}

func RequestLog(t *testing.T) option.RequestOption {
	t.Helper()
	return shared.RequestLog(t)
}

func PollPause(ctx context.Context) error { return shared.PollPause(ctx) }
