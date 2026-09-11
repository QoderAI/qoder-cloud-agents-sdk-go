package testutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
)

// SafeError omits response bodies and URLs: either may contain a token or a
// presigned storage grant. Request IDs are enough to correlate server diagnostics.
func SafeError(err error) string {
	if errors.Is(err, context.Canceled) {
		return "request canceled"
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "request deadline exceeded"
	}
	var api *apierror.Error
	if errors.As(err, &api) {
		return fmt.Sprintf("API failure: status=%d type=%s code=%s request_id=%s", api.StatusCode, api.Type(), api.Code, api.RequestID)
	}
	return fmt.Sprintf("request failed (%T); response body and URL omitted", err)
}

// CleanupFailure identifies an incomplete cleanup step, not an already deleted
// target resource (for example, a batch output which has not been generated).
type CleanupFailure struct{ Err error }

func (e *CleanupFailure) Error() string { return "cleanup incomplete: " + SafeError(e.Err) }
func (e *CleanupFailure) Unwrap() error { return e.Err }
func ResourceAlreadyGone(err error) bool {
	var incomplete *CleanupFailure
	if errors.As(err, &incomplete) {
		return false
	}
	var api *apierror.Error
	return errors.As(err, &api) && api.StatusCode == 404
}
