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

// SkillVersionService provides Forward SkillVersion operations.
type SkillVersionService struct {
	Options []option.RequestOption
}

func NewSkillVersionService(opts ...option.RequestOption) SkillVersionService {
	return SkillVersionService{Options: slices.Clone(opts)}
}

// List Skill versions
func (r *SkillVersionService) List(ctx context.Context, id string, params SkillVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[SkillVersion], err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("skills/%s/versions", url.PathEscape(id))
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
func (r *SkillVersionService) ListAutoPaging(ctx context.Context, id string, params SkillVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[SkillVersion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, id, params, opts...))
}

type SkillVersionListParams struct {
	// Page size, maximum 100, default 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for paging forward; take the value from `next_page` in the previous
	// response. Omit to start from the first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	paramObj
}

func (r SkillVersionListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Skill version
func (r *SkillVersionService) New(ctx context.Context, id string, params SkillVersionNewParams, opts ...option.RequestOption) (res *SkillVersion, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("skills/%s/versions", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type SkillVersionNewParams struct {
	// Upload field, which may appear **multiple** times. Two shapes are supported: • a
	// single `.zip` archive; • a bare file tree, where every part uploads one file and
	// `filename` carries its relative path (e.g. `customer-reply/SKILL.md`,
	// `customer-reply/scripts/run.sh`). Neither the archive itself nor the total
	// uncompressed size may exceed 50 MB.
	Files []io.Reader `json:"files" api:"required" format:"binary"`
	paramObj
}

func (r SkillVersionNewParams) MarshalMultipart() ([]byte, string, error) {
	return marshalMultipart(r, r.ExtraFields())
}

// Get Skill version
func (r *SkillVersionService) Get(ctx context.Context, id string, version string, opts ...option.RequestOption) (res *SkillVersion, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}
	if version == "" {
		return nil, fmt.Errorf("missing required version parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("skills/%s/versions/%s", url.PathEscape(id), url.PathEscape(version))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Delete Skill version
func (r *SkillVersionService) Delete(ctx context.Context, id string, version string, opts ...option.RequestOption) (res *DeletedSkillVersion, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}
	if version == "" {
		return nil, fmt.Errorf("missing required version parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("skills/%s/versions/%s", url.PathEscape(id), url.PathEscape(version))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Download Skill version content
func (r *SkillVersionService) Download(ctx context.Context, id string, version string, opts ...option.RequestOption) (res *http.Response, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}
	if version == "" {
		return nil, fmt.Errorf("missing required version parameter")
	}

	opts = slices.Concat(r.Options, []option.RequestOption{option.WithHeader("Accept", "application/zip")}, opts)
	path := fmt.Sprintf("skills/%s/versions/%s/content", url.PathEscape(id), url.PathEscape(version))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type DeletedSkillVersion struct {
	// The deleted version number.
	ID string `json:"id"`
	// Always `"skill_version_deleted"`
	Type string `json:"type"`
	// Always `true`
	Deleted bool `json:"deleted"`
	JSON    struct {
		ID          respjson.Field
		Type        respjson.Field
		Deleted     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DeletedSkillVersion) RawJSON() string                  { return r.JSON.raw }
func (r *DeletedSkillVersion) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SkillVersion struct {
	ID            string         `json:"id"`
	SkillID       string         `json:"skill_id"`
	Version       string         `json:"version"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Directory     string         `json:"directory"`
	ContentSize   int64          `json:"content_size"`
	ContentSHA256 string         `json:"content_sha256"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at" format:"date-time"`
	Metadata      map[string]any `json:"metadata"`
	JSON          struct {
		ID            respjson.Field
		SkillID       respjson.Field
		Version       respjson.Field
		Name          respjson.Field
		Description   respjson.Field
		Directory     respjson.Field
		ContentSize   respjson.Field
		ContentSHA256 respjson.Field
		Status        respjson.Field
		CreatedAt     respjson.Field
		Metadata      respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r SkillVersion) RawJSON() string                  { return r.JSON.raw }
func (r *SkillVersion) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
