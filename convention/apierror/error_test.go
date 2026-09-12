package apierror_test

import (
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
)

func TestErrorUnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		name                            string
		body                            string
		errorType, message, code, reqID string
	}{
		{
			name:      "nested_service_envelope",
			body:      `{"type":"error","request_id":"request-from-body","error":{"type":"not_found_error","message":"Agent 'ag-1' was not found.","code":"NOT_FOUND"}}`,
			errorType: "not_found_error", message: "Agent 'ag-1' was not found.", code: "NOT_FOUND", reqID: "request-from-body",
		},
		{
			// The gateway rejects unauthenticated calls before they reach the
			// service, reporting the code at the top level and no request_id.
			name: "flat_gateway_envelope",
			body: `{"code":"TOKEN_INVALID","message":"invalid pt-token","timestamp":"1789199957397"}`,
			code: "TOKEN_INVALID", message: "invalid pt-token", reqID: "request-from-header",
		},
		{
			// Top-level "type" is the envelope discriminator, never an error type.
			name:      "nested_envelope_wins_over_top_level",
			body:      `{"type":"error","message":"outer","code":"OUTER","error":{"type":"invalid_request_error","message":"inner","code":"INNER"}}`,
			errorType: "invalid_request_error", message: "inner", code: "INNER", reqID: "request-from-header",
		},
		{
			name:    "non_json_body",
			body:    "upstream connect error",
			message: "upstream connect error", reqID: "request-from-header",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := apierror.Error{StatusCode: 401, RequestID: "request-from-header"}
			if err := e.UnmarshalJSON([]byte(tc.body)); err != nil {
				t.Fatal(err)
			}
			if e.Type() != tc.errorType || e.Message != tc.message || e.Code != tc.code || e.RequestID != tc.reqID || e.RawJSON() != tc.body {
				t.Fatal(e.Type(), e.Message, e.Code, e.RequestID, e.RawJSON())
			}
		})
	}
}
