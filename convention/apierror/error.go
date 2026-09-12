// Package apierror describes Qoder's standard error envelope.
package apierror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error is an HTTP API error. Transport and context errors remain unwrapped.
type Error struct {
	StatusCode int            `json:"-"`
	Request    *http.Request  `json:"-"`
	Response   *http.Response `json:"-"`
	RequestID  string         `json:"request_id"`
	Message    string         `json:"-"`
	Code       string         `json:"-"`
	errorType  string
	raw        string
}

func (e *Error) Type() string   { return e.errorType }
func (e Error) RawJSON() string { return e.raw }
func (e *Error) UnmarshalJSON(data []byte) error {
	e.raw = string(data)
	// The service nests its error under "error", but the auth gateway rejects
	// requests before they reach the service using a flat body. The nested
	// envelope's top-level "type" is the constant discriminator "error", so it
	// is never an error type and is not read here.
	var body struct {
		RequestID string `json:"request_id"`
		Message   string `json:"message"`
		Code      string `json:"code"`
		Error     struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(data, &body); err != nil {
		e.Message = string(data)
		return nil
	}
	if body.RequestID != "" {
		e.RequestID = body.RequestID
	}
	e.errorType = body.Error.Type
	e.Message = body.Error.Message
	e.Code = body.Error.Code
	if e.Message == "" {
		e.Message = body.Message
	}
	if e.Code == "" {
		e.Code = body.Code
	}
	return nil
}
func (e *Error) Error() string {
	return fmt.Sprintf("qoder: HTTP %d %s: %s (request_id=%s)", e.StatusCode, e.errorType, e.Message, e.RequestID)
}
