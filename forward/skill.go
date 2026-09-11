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

// SkillService provides Forward Skill operations.
type SkillService struct {
	Options  []option.RequestOption
	Versions SkillVersionService
}

func NewSkillService(opts ...option.RequestOption) SkillService {
	return SkillService{Options: slices.Clone(opts), Versions: NewSkillVersionService(opts...)}
}

// 列出 Skill.
func (r *SkillService) List(ctx context.Context, params SkillListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Skill], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "skills"
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
func (r *SkillService) ListAutoPaging(ctx context.Context, params SkillListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Skill] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

type SkillListParams struct {
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// 向后翻页游标；与 `page`、`before_id` 互斥。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标；与 `page`、`after_id` 互斥。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 按 Skill 展示名前缀搜索，不区分大小写。
	DisplayTitle param.Opt[string] `query:"display_title,omitzero" json:"-"`
	// 按 Skill 来源过滤，可选 `custom`、`qoder`。传 `source` 时不支持 `before_id`。
	Source param.Opt[string] `query:"source,omitzero" json:"-"`
	// ⚠️ **已弃用**：`display_title` 的兼容别名，语义完全一致。请使用 `display_title`。
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	paramObj
}

func (r SkillListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Skill.
func (r *SkillService) New(ctx context.Context, params SkillNewParams, opts ...option.RequestOption) (res *Skill, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "skills"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type SkillNewParams struct {
	// 推荐上传字段，可**重复出现**多次。支持两种形态： ① 单个 `.zip` 包； ② 裸文件树——每个 part 独立上传一个文件，`filename` 携带相对路径（如 `code-review/SKILL.md`、`code-review/scripts/run.sh`）。 压缩包本身与解压后总大小均不超过 50 MB。
	Files []io.Reader `json:"files,omitzero" format:"binary"`
	// 调用方元数据对象，最多 15 个键；`created_by` 为保留字段，不可传入（传入返回 400）。
	Metadata map[string]any `json:"metadata,omitzero" api:"metadata"`
	// Forward Resource icon 公开 ID。
	IconID param.Opt[string] `json:"icon_id,omitzero"`
	// ⚠️ **已弃用**：单个 `.zip` 包，宽松包规则。命中时响应头返回 `Deprecation: true`。请迁移到 `files`。
	File io.Reader `json:"file,omitzero" format:"binary"`
	// ⚠️ **已弃用**：最终名称始终从上传包内 `SKILL.md` frontmatter 的 `name` 解析。字段保留仅为兼容，传入将被忽略。
	Name param.Opt[string] `json:"name,omitzero"`
	// ⚠️ **已弃用**：最终描述始终从 `SKILL.md` 解析。
	Description param.Opt[string] `json:"description,omitzero"`
	// ⚠️ **已弃用**：Skill 创建类型，可选 `custom`、`prebuilt`，默认 `custom`。`prebuilt` 会使响应 `source` 字段返回 `qoder`（其余为 `custom`）。
	Type param.Opt[string] `json:"type,omitzero"`
	// 建议提供。相同 key 且规范化后的 `files` 指纹一致时可安全重试。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SkillNewParams) MarshalMultipart() ([]byte, string, error) {
	return marshalMultipart(r, r.ExtraFields())
}

// 查询 Skill.
func (r *SkillService) Get(ctx context.Context, id string, params SkillGetParams, opts ...option.RequestOption) (res *Skill, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("skills/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

type SkillGetParams struct {
	// ⚠️ **已弃用**：为 `true` 时随响应返回 `content` 与 `content_encoding`（base64 zip）。命中时响应头会返回 `Deprecation: true`。请改用 [下载 Skill 版本内容](./Versions/download.md)。
	IncludeContent param.Opt[bool] `query:"include_content,omitzero" json:"-"`
	paramObj
}

func (r SkillGetParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 修改 Skill.
func (r *SkillService) Update(ctx context.Context, id string, params SkillUpdateParams, opts ...option.RequestOption) (res *Skill, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("skills/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodPut, path, params, &res, opts...)
	return res, err
}

type SkillUpdateParams struct {
	// 新描述。
	Description param.Opt[string] `json:"description,omitzero"`
	// 新内容（zip 包内容）。压缩包本身与解压后总大小均不超过 50 MB，超过返回 400；请求体整体（含 base64 编码与 JSON 信封）上限约 67.7 MB，超过返回 413。
	Content param.Opt[string] `json:"content,omitzero"`
	// `content` 的编码。支持 `base64`、`utf-8`、`utf8`、`plain`、`text`；省略时按 UTF-8 文本处理。传入该字段时必须同时提供非空 `content`。
	ContentEncoding param.Opt[string] `json:"content_encoding,omitzero"`
	// 元数据对象，会**替换**当前 metadata（非合并）；传入时不能为 `null`，value 必须为 string。`created_by` 为保留字段，不可传入（传入返回 400）。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 更新或清空 Forward icon。
	IconID param.Opt[string] `json:"icon_id,omitzero" api:"nullable"`
	// ⚠️ **已弃用**：技能名不可修改。传入必须与当前规范名完全一致，否则返回 400；一致时为空操作。
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r SkillUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow SkillUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SkillUpdateParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// 删除 Skill.
func (r *SkillService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	if id == "" {
		return fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("skills/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type Skill struct {
	ID            string           `json:"id"`
	Type          string           `json:"type"`
	DisplayTitle  string           `json:"display_title"`
	Description   string           `json:"description"`
	Source        string           `json:"source"`
	LatestVersion string           `json:"latest_version"`
	Metadata      map[string]any   `json:"metadata"`
	CreatedAt     time.Time        `json:"created_at" format:"date-time"`
	UpdatedAt     time.Time        `json:"updated_at" format:"date-time"`
	IdentityID    string           `json:"identity_id" api:"nullable"`
	IconURL       string           `json:"icon_url" api:"nullable"`
	BindingInfo   map[string]int64 `json:"binding_info"`
	JSON          struct {
		ID            respjson.Field
		Type          respjson.Field
		DisplayTitle  respjson.Field
		Description   respjson.Field
		Source        respjson.Field
		LatestVersion respjson.Field
		Metadata      respjson.Field
		CreatedAt     respjson.Field
		UpdatedAt     respjson.Field
		IdentityID    respjson.Field
		IconURL       respjson.Field
		BindingInfo   respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r Skill) RawJSON() string                  { return r.JSON.raw }
func (r *Skill) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
