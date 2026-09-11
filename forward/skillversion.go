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

// 列出 Skill 版本.
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
	// 分页大小，最大 100，默认 20。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向后翻页游标；取值来自上一页响应的 `next_page`；不传即从第一页开始。
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	paramObj
}

func (r SkillVersionListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Skill 版本.
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
	// 上传字段，可**重复**出现多次。支持两种形态： • 单个 `.zip` 包； • 裸文件树——每个 part 独立上传一个文件，`filename` 携带相对路径（如 `customer-reply/SKILL.md`、`customer-reply/scripts/run.sh`）。 压缩包本身与解压后总大小均不超过 50 MB。
	Files []io.Reader `json:"files" api:"required" format:"binary"`
	paramObj
}

func (r SkillVersionNewParams) MarshalMultipart() ([]byte, string, error) {
	return marshalMultipart(r, r.ExtraFields())
}

// 查询 Skill 版本.
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

// 删除 Skill 版本.
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

// 下载 Skill 版本内容.
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
	// 被删除的版本号。
	ID string `json:"id"`
	// 固定值 `"skill_version_deleted"`
	Type string `json:"type"`
	// 固定值 `true`
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
