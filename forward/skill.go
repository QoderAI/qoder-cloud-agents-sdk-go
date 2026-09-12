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

// List Skills
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
	// Page size, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Pagination cursor (recommended), taken from `next_page` in the previous response;
	// mutually exclusive with `after_id` and `before_id`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Cursor for the next page; mutually exclusive with `page` and `before_id`.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page; mutually exclusive with `page` and `after_id`.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Search by Skill display title prefix, case-insensitive.
	DisplayTitle param.Opt[string] `query:"display_title,omitzero" json:"-"`
	// Filter by Skill source, one of `custom`, `qoder`. `before_id` is not supported
	// when `source` is set.
	Source param.Opt[string] `query:"source,omitzero" json:"-"`
	// ⚠️ **Deprecated**: compatibility alias for `display_title` with identical
	// semantics. Use `display_title` instead.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	paramObj
}

func (r SkillListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Skill
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
	// Recommended upload field, which may **repeat** several times. Two shapes are
	// supported: (1) a single `.zip` archive; (2) a bare file tree — each part uploads
	// one file and carries its relative path in `filename` (such as
	// `code-review/SKILL.md`, `code-review/scripts/run.sh`). Neither the archive itself
	// nor its uncompressed total size may exceed 50 MB.
	Files []io.Reader `json:"files,omitzero" format:"binary"`
	// Caller metadata object, at most 15 keys. `created_by` is reserved and must not be
	// sent (sending it returns 400).
	Metadata map[string]any `json:"metadata,omitzero" api:"metadata"`
	// Public ID of the Forward Resource icon.
	IconID param.Opt[string] `json:"icon_id,omitzero"`
	// ⚠️ **Deprecated**: a single `.zip` archive with relaxed packaging rules. Responses
	// carry the `Deprecation: true` header when it is used. Migrate to `files`.
	File io.Reader `json:"file,omitzero" format:"binary"`
	// ⚠️ **Deprecated**: the final name is always parsed from the `name` in the
	// `SKILL.md` frontmatter of the uploaded package. The field is kept for
	// compatibility only and is ignored when sent.
	Name param.Opt[string] `json:"name,omitzero"`
	// ⚠️ **Deprecated**: the final description is always parsed from `SKILL.md`.
	Description param.Opt[string] `json:"description,omitzero"`
	// ⚠️ **Deprecated**: Skill creation type, one of `custom`, `prebuilt`, defaults to
	// `custom`. `prebuilt` makes the response `source` field return `qoder` (everything
	// else returns `custom`).
	Type param.Opt[string] `json:"type,omitzero"`
	// Recommended. Retrying is safe when the same key is paired with an identical
	// normalized `files` fingerprint.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SkillNewParams) MarshalMultipart() ([]byte, string, error) {
	return marshalMultipart(r, r.ExtraFields())
}

// Get Skill
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
	// ⚠️ **Deprecated**: when `true`, the response also carries `content` and
	// `content_encoding` (a base64 zip). Responses carry the `Deprecation: true` header
	// when it is used. Use [Download Skill version content](./Versions/download.md)
	// instead.
	IncludeContent param.Opt[bool] `query:"include_content,omitzero" json:"-"`
	paramObj
}

func (r SkillGetParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Update Skill
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
	// New description.
	Description param.Opt[string] `json:"description,omitzero"`
	// New content (the contents of a zip archive). Neither the archive itself nor its
	// uncompressed total size may exceed 50 MB, otherwise 400 is returned. The request
	// body as a whole (including base64 encoding and the JSON envelope) is capped at
	// about 67.7 MB, beyond which 413 is returned.
	Content param.Opt[string] `json:"content,omitzero"`
	// Encoding of `content`. Supports `base64`, `utf-8`, `utf8`, `plain`, `text`; when
	// omitted the value is treated as UTF-8 text. Sending this field requires a
	// non-empty `content` as well.
	ContentEncoding param.Opt[string] `json:"content_encoding,omitzero"`
	// Metadata object that **replaces** the current metadata (not a merge). It must not
	// be `null` when sent, and each value must be a string. `created_by` is reserved and
	// must not be sent (sending it returns 400).
	Metadata map[string]any `json:"metadata,omitzero"`
	// Update or clear the Forward icon.
	IconID param.Opt[string] `json:"icon_id,omitzero" api:"nullable"`
	// ⚠️ **Deprecated**: the Skill name cannot be changed. A value sent must match the
	// current canonical name exactly, otherwise 400 is returned; a matching value is a
	// no-op.
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r SkillUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow SkillUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SkillUpdateParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// Delete Skill
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
