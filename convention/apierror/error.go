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
	// WorkspaceID is retained for source compatibility with the upstream error shape.
	WorkspaceID string `json:"-"`
	Message     string `json:"-"`
	Code        string `json:"-"`
	errorType   string
	raw         string
}

func (e *Error) Type() string   { return e.errorType }
func (e Error) RawJSON() string { return e.raw }
func (e *Error) UnmarshalJSON(data []byte) error {
	e.raw = string(data)
	var body struct {
		RequestID string `json:"request_id"`
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
	return nil
}
func (e *Error) Error() string {
	return fmt.Sprintf("qoder: HTTP %d %s: %s (request_id=%s)", e.StatusCode, e.errorType, e.Message, e.RequestID)
}
