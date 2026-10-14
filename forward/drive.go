package forward

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// DriveService manages persistent Identity files. Drive is a Beta API.
// Every method requires IdentityID, including when using an Identity SAT.
type DriveService struct{ Options []option.RequestOption }

func NewDriveService(opts ...option.RequestOption) DriveService {
	return DriveService{Options: slices.Clone(opts)}
}

// ListEntries returns one directory page in service order. Pass NextPageToken
// unchanged as PageToken to continue within the same Identity and directory.
func (r *DriveService) ListEntries(ctx context.Context, params DriveListEntriesParams, opts ...option.RequestOption) (res *DriveListEntriesResponse, err error) {
	cfg, err := convention.NewRequestConfig(ctx, http.MethodGet, "drives/entries", params, &res, slices.Concat(r.Options, opts)...)
	if err != nil {
		return nil, err
	}
	// Drive requires literal slashes in the path query, after request options
	// have been applied. Other values (especially cursors) stay encoded.
	fields := strings.Split(cfg.Request.URL.RawQuery, "&")
	for i, field := range fields {
		key, value, _ := strings.Cut(field, "=")
		if key != "path" {
			continue
		}
		decoded, decodeErr := url.QueryUnescape(value)
		if decodeErr != nil {
			return nil, decodeErr
		}
		fields[i] = "path=" + strings.ReplaceAll(url.QueryEscape(decoded), "%2F", "/")
	}
	cfg.Request.URL.RawQuery = strings.Join(fields, "&")
	err = cfg.Execute()
	return res, err
}

// CreateUploadURL returns a short-lived PUT URL. Use only its returned headers for the storage request. Replaying an Idempotency-Key does not renew expiry.
func (r *DriveService) CreateUploadURL(ctx context.Context, params DriveCreateUploadURLParams, opts ...option.RequestOption) (res *DriveUploadURL, err error) {
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, "drives/upload-url", params, &res, slices.Concat(r.Options, opts)...)
	return res, err
}

// CreateDownloadURL returns a short-lived GET URL. Use a new Idempotency-Key or omit it to refresh an expired URL.
func (r *DriveService) CreateDownloadURL(ctx context.Context, params DriveCreateDownloadURLParams, opts ...option.RequestOption) (res *DriveDownloadURL, err error) {
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, "drives/download-url", params, &res, slices.Concat(r.Options, opts)...)
	return res, err
}

// DeleteEntry deletes a file or directory. Recursive must be true for a nonempty directory. Partial failures retain deleted_count in the API error raw body.
func (r *DriveService) DeleteEntry(ctx context.Context, params DriveDeleteEntryParams, opts ...option.RequestOption) (res *DriveDeleteEntryResponse, err error) {
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, "drives/entries", params, &res, slices.Concat(r.Options, opts)...)
	return res, err
}

// Clear clears the selected Identity Drive. Use only with an Identity whose files you intend to delete.
func (r *DriveService) Clear(ctx context.Context, params DriveClearParams, opts ...option.RequestOption) (res *DriveClearResponse, err error) {
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, "drives", params, &res, slices.Concat(r.Options, opts)...)
	return res, err
}

type DriveListEntriesParams struct {
	IdentityID string            `query:"identity_id" json:"-" api:"required"`
	Path       param.Opt[string] `query:"path,omitzero" json:"-"`
	Limit      param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	PageToken  param.Opt[string] `query:"page_token,omitzero" json:"-"`
	paramObj
}

func (r *DriveListEntriesParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r DriveListEntriesParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type DriveCreateUploadURLParams struct {
	IdentityID  string            `query:"identity_id" json:"-" api:"required"`
	Path        string            `json:"path" api:"required"`
	ContentType param.Opt[string] `json:"content_type,omitzero"`
	paramObj
}

func (r DriveCreateUploadURLParams) MarshalJSON() ([]byte, error) {
	type shadow DriveCreateUploadURLParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DriveCreateUploadURLParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r DriveCreateUploadURLParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type DriveCreateDownloadURLParams struct {
	IdentityID string `query:"identity_id" json:"-" api:"required"`
	Path       string `json:"path" api:"required"`
	paramObj
}

func (r DriveCreateDownloadURLParams) MarshalJSON() ([]byte, error) {
	type shadow DriveCreateDownloadURLParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DriveCreateDownloadURLParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r DriveCreateDownloadURLParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type DriveDeleteEntryParams struct {
	IdentityID string          `query:"identity_id" json:"-" api:"required"`
	Path       string          `json:"path" api:"required"`
	Recursive  param.Opt[bool] `json:"recursive,omitzero"`
	paramObj
}

func (r DriveDeleteEntryParams) MarshalJSON() ([]byte, error) {
	type shadow DriveDeleteEntryParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DriveDeleteEntryParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r DriveDeleteEntryParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type DriveClearParams struct {
	IdentityID string `query:"identity_id" json:"-" api:"required"`
	paramObj
}

func (r *DriveClearParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
func (r DriveClearParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// DriveEntry represents a file or directory. JSON.Size.Valid() distinguishes a zero-byte file from a directory.
type DriveEntry struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Size         int64     `json:"size"`
	ETag         string    `json:"etag"`
	LastModified time.Time `json:"last_modified" format:"date-time"`
	JSON         struct {
		Path         respjson.Field
		Name         respjson.Field
		Type         respjson.Field
		Size         respjson.Field
		ETag         respjson.Field
		LastModified respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r DriveEntry) RawJSON() string                  { return r.JSON.raw }
func (r *DriveEntry) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DriveListEntriesResponse struct {
	Entries       []DriveEntry `json:"entries"`
	NextPageToken string       `json:"next_page_token"`
	JSON          struct {
		Entries       respjson.Field
		NextPageToken respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r DriveListEntriesResponse) RawJSON() string { return r.JSON.raw }
func (r *DriveListEntriesResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DriveUploadURL struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ExpiresAt time.Time         `json:"expires_at" format:"date-time"`
	JSON      struct {
		URL         respjson.Field
		Method      respjson.Field
		Headers     respjson.Field
		ExpiresAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DriveUploadURL) RawJSON() string                  { return r.JSON.raw }
func (r *DriveUploadURL) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DriveDownloadURL struct {
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers"`
	ExpiresAt   time.Time         `json:"expires_at" format:"date-time"`
	Size        int64             `json:"size"`
	ETag        string            `json:"etag"`
	ContentType string            `json:"content_type"`
	JSON        struct {
		URL         respjson.Field
		Method      respjson.Field
		Headers     respjson.Field
		ExpiresAt   respjson.Field
		Size        respjson.Field
		ETag        respjson.Field
		ContentType respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DriveDownloadURL) RawJSON() string                  { return r.JSON.raw }
func (r *DriveDownloadURL) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DriveDeleteEntryResponse struct {
	DeletedCount int64 `json:"deleted_count"`
	JSON         struct {
		DeletedCount respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r DriveDeleteEntryResponse) RawJSON() string { return r.JSON.raw }
func (r *DriveDeleteEntryResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DriveClearResponse struct {
	DeletedCount int64 `json:"deleted_count"`
	JSON         struct {
		DeletedCount respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r DriveClearResponse) RawJSON() string                  { return r.JSON.raw }
func (r *DriveClearResponse) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
