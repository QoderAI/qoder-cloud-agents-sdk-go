package pagination

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	requestconfig "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
)

// NewSearchPage sends a resource search with the required Beta header and JSON body pagination.
func NewSearchPage[T any](ctx context.Context, path string, params any, opts ...option.RequestOption) (res *PageCursor[T], err error) {
	var raw *http.Response
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	const beta = "search-2026-08-31"
	values := strings.Join(cfg.Request.Header.Values("x-qoder-beta"), ",")
	found := false
	for _, value := range strings.Split(values, ",") {
		if strings.TrimSpace(value) == beta {
			found = true
		}
	}
	if !found {
		if values != "" {
			values += ","
		}
		values += beta
	}
	cfg.Request.Header.Set("x-qoder-beta", values)
	if err = cfg.Execute(); err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	res.bodyPagination = true
	return res, nil
}

func advanceSearchBody(cfg *requestconfig.RequestConfig, next string) error {
	if cfg.Request.Body == nil {
		return fmt.Errorf("search pagination request has no body")
	}
	data, err := io.ReadAll(cfg.Request.Body)
	_ = cfg.Request.Body.Close()
	if err != nil {
		return err
	}
	var body map[string]json.RawMessage
	if err = json.Unmarshal(data, &body); err != nil {
		return err
	}
	if body == nil {
		return fmt.Errorf("search pagination body must be a JSON object")
	}
	var current string
	_ = json.Unmarshal(body["page"], &current)
	if current == next {
		return fmt.Errorf("pagination cursor did not advance")
	}
	body["page"], err = json.Marshal(next)
	if err != nil {
		return err
	}
	data, err = json.Marshal(body)
	if err != nil {
		return err
	}
	cfg.Request.Body = nil
	cfg.Request.GetBody = nil
	cfg.Body = bytes.NewBuffer(data)
	return nil
}
