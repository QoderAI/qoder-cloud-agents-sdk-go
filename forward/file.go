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

// 列出 File.
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
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// 向后翻页游标；与 `page`、`before_id` 互斥。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标；与 `page`、`after_id` 互斥。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 按文件名搜索。
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	// 按资源作用域 ID 过滤，常用于 Session 资源文件查询。传入时不要同时使用 `before_id` 或 `after_id`；当前游标参数在该过滤模式下不生效。
	ScopeID param.Opt[string] `query:"scope_id,omitzero" json:"-"`
	paramObj
}

func (r FileListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 上传 File.
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
	// 待上传文件内容。支持类型见[支持上传的文件类型](./schemas.md#支持上传的文件类型)。
	File io.Reader `json:"file" api:"required" format:"binary"`
	// 文件展示名，未传时使用 multipart 文件名；规范化后长度为 1-255 bytes。
	Name param.Opt[string] `json:"name,omitzero"`
	// 文件用途，默认 `user_upload`；作为 Batch 输入文件时必须传 `session_resource`。
	Purpose param.Opt[string] `json:"purpose,omitzero"`
	// 元数据对象；`created_by` 为保留字段，不可传入（传入返回 400）。
	Metadata map[string]any `json:"metadata,omitzero" api:"metadata"`
	// 可选创建请求幂等键。传入时相同 key 只能用于相同请求；不传时不提供本地幂等重放保护。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r FileUploadParams) MarshalMultipart() ([]byte, string, error) {
	return marshalMultipart(r, r.ExtraFields())
}

// 查询 File.
func (r *FileService) GetMetadata(ctx context.Context, fileID string, opts ...option.RequestOption) (res *FileMetadata, err error) {
	if fileID == "" {
		return nil, fmt.Errorf("missing required file_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("files/%s", url.PathEscape(fileID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 删除 File.
func (r *FileService) Delete(ctx context.Context, fileID string, opts ...option.RequestOption) (err error) {
	if fileID == "" {
		return fmt.Errorf("missing required file_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("files/%s", url.PathEscape(fileID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// 下载 File.
func (r *FileService) Download(ctx context.Context, fileID string, opts ...option.RequestOption) (res *http.Response, err error) {
	if fileID == "" {
		return nil, fmt.Errorf("missing required file_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("files/%s/content", url.PathEscape(fileID))
	return convention.DownloadFile(ctx, path, opts...)
}

type FileMetadata struct {
	// File ID。
	ID string `json:"id"`
	// 固定为 `file`。
	Type string `json:"type"`
	// 文件名。
	Filename string `json:"filename"`
	// 文件大小，单位为 byte。
	SizeBytes int64 `json:"size_bytes"`
	// MIME 类型。
	MIMEType string `json:"mime_type"`
	// 创建时间，RFC 3339 格式。
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// 最后更新时间，RFC 3339 格式。
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	// 是否可下载。
	Downloadable bool `json:"downloadable"`
	// 文件关联的资源作用域，如 Session。
	Scope map[string]any `json:"scope" api:"nullable"`
	// 文件元数据。
	Metadata map[string]any `json:"metadata"`
	// Forward 归属身份。
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
