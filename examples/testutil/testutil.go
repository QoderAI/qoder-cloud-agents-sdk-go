// Package testutil exposes shared test helpers to SDK tests outside examples,
// which cannot import examples/internal/testutil directly.
package testutil

import (
	"encoding/json"
	"reflect"
	"testing"

	shared "github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/testutil"
)

type Endpoint = shared.Endpoint
type Transport = shared.Transport
type TurnResult = shared.TurnResult
type CleanupFailure = shared.CleanupFailure

func FailureContracts(t *testing.T, endpoints []Endpoint, client func(Transport) any) {
	t.Helper()
	shared.FailureContracts(t, endpoints, client)
}

func CheckFields(t *testing.T, value reflect.Value, path string) {
	t.Helper()
	shared.CheckFields(t, value, path)
}

func InvokeJSON(t *testing.T, client any, endpoint Endpoint, body json.RawMessage) (any, error) {
	t.Helper()
	return shared.InvokeJSON(t, client, endpoint, body)
}

func SafeError(err error) string         { return shared.SafeError(err) }
func ResourceAlreadyGone(err error) bool { return shared.ResourceAlreadyGone(err) }
