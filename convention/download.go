package convention

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// FileContentLink is Qoder's temporary download grant.
type FileContentLink struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// DownloadFile obtains a temporary link and returns its streaming HTTP response.
// The caller closes Body. API authorization and middleware are never forwarded
// to the storage host; the signed URL provides the storage authorization.
func DownloadFile(ctx context.Context, path string, opts ...RequestOption) (*http.Response, error) {
	var link FileContentLink
	cfg, err := NewRequestConfig(ctx, http.MethodGet, path, nil, &link, opts...)
	if err != nil {
		return nil, err
	}
	if err = cfg.Execute(); err != nil {
		return nil, err
	}
	u, err := url.Parse(link.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid file content link: %w", err)
	}
	if (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil {
		return nil, fmt.Errorf("invalid file content link")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var res *http.Response
	download := RequestConfig{Context: ctx, Request: req, BaseURL: u, HTTPClient: cfg.HTTPClient, CustomHTTPDoer: cfg.CustomHTTPDoer, RequestTimeout: cfg.RequestTimeout, MaxRetries: cfg.MaxRetries, ResponseBodyInto: &res, ResponseInto: cfg.ResponseInto}
	err = download.Execute()
	return res, err
}
