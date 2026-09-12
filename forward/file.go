package forward

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// FileService provides Forward File operations.
type FileService struct {
	Options []option.RequestOption
}

func NewFileService(opts ...option.RequestOption) FileService {
	return FileService{Options: slices.Clone(opts)}
}

// List Files
func (r *FileService) List(ctx context.Context, params FileListParams, opts ...option.RequestOption) (res *pagination.PageCursor[FileMetadata], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "files"
	var raw *http.Response
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	cfg, err := convention.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	if err = cfg.Execute(); err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}
func (r *FileService) ListAutoPaging(ctx context.Context, params FileListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[FileMetadata] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

type FileListParams struct {
	// Page size, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Pagination cursor (recommended), taken from `next_page` in the previous response;
	// mutually exclusive with `after_id` and `before_id`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Cursor for the next page; mutually exclusive with `page` and `before_id`.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page; mutually exclusive with `page` and `after_id`.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Search by filename.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	// Filter by resource scope ID, commonly used to look up Session resource files. Do
	// not combine it with `before_id` or `after_id`; cursor parameters currently have no
	// effect in this filter mode.
	ScopeID param.Opt[string] `query:"scope_id,omitzero" json:"-"`
	paramObj
}

func (r FileListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Upload File
func (r *FileService) Upload(ctx context.Context, params FileUploadParams, opts ...option.RequestOption) (res *FileMetadata, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "files"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type FileUploadParams struct {
	// Content of the file to upload. For accepted types see
	// [Supported upload file types](./schemas.md).
	File io.Reader `json:"file" api:"required" format:"binary"`
	// Display name for the file; the multipart filename is used when omitted. After
	// normalization the length must be 1-255 bytes.
	Name param.Opt[string] `json:"name,omitzero"`
	// Purpose of the file, defaults to `user_upload`; must be `session_resource` when the
	// file is used as Batch input.
	Purpose param.Opt[string] `json:"purpose,omitzero"`
	// Metadata object. `created_by` is reserved and must not be sent (sending it returns
	// 400).
	Metadata map[string]any `json:"metadata,omitzero" api:"metadata"`
	// Optional idempotency key for the create request. When set, the same key may only be
	// reused for an identical request; when omitted no local idempotent replay protection
	// is provided.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r FileUploadParams) MarshalMultipart() ([]byte, string, error) {
	return marshalMultipart(r, r.ExtraFields())
}

// Get File
func (r *FileService) GetMetadata(ctx context.Context, fileID string, opts ...option.RequestOption) (res *FileMetadata, err error) {
	if fileID == "" {
		return nil, fmt.Errorf("missing required file_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("files/%s", url.PathEscape(fileID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Delete File
func (r *FileService) Delete(ctx context.Context, fileID string, opts ...option.RequestOption) (err error) {
	if fileID == "" {
		return fmt.Errorf("missing required file_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("files/%s", url.PathEscape(fileID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Download File
func (r *FileService) Download(ctx context.Context, fileID string, opts ...option.RequestOption) (res *http.Response, err error) {
	if fileID == "" {
		return nil, fmt.Errorf("missing required file_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("files/%s/content", url.PathEscape(fileID))
	return convention.DownloadFile(ctx, path, opts...)
}

type FileMetadata struct {
	// File ID.
	ID string `json:"id"`
	// Always `file`.
	Type string `json:"type"`
	// Filename.
	Filename string `json:"filename"`
	// File size in bytes.
	SizeBytes int64 `json:"size_bytes"`
	// MIME type.
	MIMEType string `json:"mime_type"`
	// Creation time in RFC 3339 format.
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// Last update time in RFC 3339 format.
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	// Whether the file can be downloaded.
	Downloadable bool `json:"downloadable"`
	// Resource scope the file belongs to, such as a Session.
	Scope map[string]any `json:"scope" api:"nullable"`
	// File metadata.
	Metadata map[string]any `json:"metadata"`
	// Owning Forward Identity.
	IdentityID string `json:"identity_id" api:"nullable"`
	JSON       struct {
		ID           respjson.Field
		Type         respjson.Field
		Filename     respjson.Field
		SizeBytes    respjson.Field
		MIMEType     respjson.Field
		CreatedAt    respjson.Field
		UpdatedAt    respjson.Field
		Downloadable respjson.Field
		Scope        respjson.Field
		Metadata     respjson.Field
		IdentityID   respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r FileMetadata) RawJSON() string                  { return r.JSON.raw }
func (r *FileMetadata) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
