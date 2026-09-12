package forward

import (
	"context"
	"fmt"
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

// EnvironmentService provides Forward Environment operations.
type EnvironmentService struct {
	Options []option.RequestOption
}

func NewEnvironmentService(opts ...option.RequestOption) EnvironmentService {
	return EnvironmentService{Options: slices.Clone(opts)}
}

// List Environments
func (r *EnvironmentService) List(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Environment], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "environments"
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
func (r *EnvironmentService) ListAutoPaging(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Environment] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

type EnvironmentListParams struct {
	// Page size, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Pagination cursor (recommended); take the value from `next_page` in the previous
	// response. Mutually exclusive with `after_id` and `before_id`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Cursor for paging forward; mutually exclusive with `page` and `before_id`.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for paging backward; mutually exclusive with `page` and `after_id`.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r EnvironmentListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Environment
func (r *EnvironmentService) New(ctx context.Context, params EnvironmentNewParams, opts ...option.RequestOption) (res *Environment, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "environments"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type EnvironmentNewParams struct {
	// Environment name; must not be empty after trimming surrounding whitespace.
	Name string `json:"name" api:"required"`
	// Description.
	Description param.Opt[string] `json:"description,omitzero"`
	// Environment runtime configuration; defaults to `{"type":"cloud"}` when omitted.
	// When passed explicitly it must not be `null` or an empty object. See the schemas
	// for the available fields.
	Config map[string]any `json:"config,omitzero"`
	// [Environment metadata](./schemas.md#environment-metadata); defaults to `{}` when
	// omitted and must not be `null` when passed explicitly.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Recommended on create requests. The same key with the same request is safe to
	// retry.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r EnvironmentNewParams) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Get Environment
func (r *EnvironmentService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Environment, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("environments/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Environment
func (r *EnvironmentService) Update(ctx context.Context, id string, params EnvironmentUpdateParams, opts ...option.RequestOption) (res *Environment, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("environments/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type EnvironmentUpdateParams struct {
	// New name.
	Name param.Opt[string] `json:"name,omitzero"`
	// New description.
	Description param.Opt[string] `json:"description,omitzero"`
	// New configuration; must not be `null` when passed, an explicit `null` returns 400.
	// See the schemas for the available fields.
	Config map[string]any `json:"config,omitzero"`
	// [Environment metadata](./schemas.md#environment-metadata) to merge; must not be
	// `null` when passed, an explicit `null` returns 400.
	Metadata map[string]any `json:"metadata,omitzero"`
	paramObj
}

func (r EnvironmentUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Archive an environment retained by historical sessions or tool calls.
func (r *EnvironmentService) Archive(ctx context.Context, id string, opts ...option.RequestOption) (res *Environment, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("environments/%s/archive", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Delete Environment
func (r *EnvironmentService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	if id == "" {
		return fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("environments/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type Environment struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Config      EnvironmentConfig `json:"config"`
	Metadata    map[string]any    `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at" format:"date-time"`
	UpdatedAt   time.Time         `json:"updated_at" format:"date-time"`
	IdentityID  string            `json:"identity_id" api:"nullable"`
	ArchivedAt  time.Time         `json:"archived_at" api:"nullable" format:"date-time"`
	JSON        struct {
		ID          respjson.Field
		Type        respjson.Field
		Name        respjson.Field
		Description respjson.Field
		Config      respjson.Field
		Metadata    respjson.Field
		CreatedAt   respjson.Field
		UpdatedAt   respjson.Field
		IdentityID  respjson.Field
		ArchivedAt  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r Environment) RawJSON() string                  { return r.JSON.raw }
func (r *Environment) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type EnvironmentConfig struct {
	Type     string                    `json:"type"`
	Packages EnvironmentConfigPackages `json:"packages"`
	JSON     struct {
		Type        respjson.Field
		Packages    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EnvironmentConfig) RawJSON() string                  { return r.JSON.raw }
func (r *EnvironmentConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type EnvironmentConfigPackages struct {
	Type  string   `json:"type"`
	Apt   []string `json:"apt"`
	Cargo []string `json:"cargo"`
	Gem   []string `json:"gem"`
	Go    []string `json:"go"`
	Npm   []string `json:"npm"`
	Pip   []string `json:"pip"`
	JSON  struct {
		Type        respjson.Field
		Apt         respjson.Field
		Cargo       respjson.Field
		Gem         respjson.Field
		Go          respjson.Field
		Npm         respjson.Field
		Pip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EnvironmentConfigPackages) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentConfigPackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
