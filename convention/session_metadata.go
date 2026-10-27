package convention

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// ExecuteNewSessionRequest adds SDK provenance to a Session creation request
// after request options have been applied. Conflicting reserved metadata fails
// before transport; caller-owned params and metadata are never modified.
func ExecuteNewSessionRequest(ctx context.Context, method string, u string, body any, dst any, opts ...RequestOption) error {
	cfg, err := NewRequestConfig(ctx, method, u, body, dst, opts...)
	if err != nil {
		return err
	}
	if cfg.Body == nil {
		return fmt.Errorf("session creation body must be a JSON object")
	}
	if closer, ok := cfg.Body.(io.Closer); ok {
		defer closer.Close()
	}
	data, err := io.ReadAll(cfg.Body)
	if err != nil {
		return err
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil || payload == nil {
		return fmt.Errorf("session creation body must be a JSON object")
	}
	metadata := map[string]json.RawMessage{}
	if raw, ok := payload["metadata"]; ok && !bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		if err := json.Unmarshal(raw, &metadata); err != nil {
			return fmt.Errorf("session creation metadata must be a JSON object or null")
		}
	}
	for _, field := range []struct{ key, value string }{
		{"__qca_sdk_created", "true"},
		{"__qca_sdk_language", "go"},
		{"__qca_sdk_version", packageVersion},
	} {
		if raw, ok := metadata[field.key]; ok {
			var value string
			if err := json.Unmarshal(raw, &value); err != nil || value != field.value {
				return fmt.Errorf("conflicting SDK metadata field: %s", field.key)
			}
		}
		metadata[field.key], _ = json.Marshal(field.value)
	}
	payload["metadata"], err = json.Marshal(metadata)
	if err != nil {
		return err
	}
	data, err = json.Marshal(payload)
	if err != nil {
		return err
	}
	cfg.Body = bytes.NewBuffer(data)
	return cfg.Execute()
}
