package forward

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// ChannelPairingService provides Forward ChannelPairing operations.
type ChannelPairingService struct {
	Options []option.RequestOption
}

func NewChannelPairingService(opts ...option.RequestOption) ChannelPairingService {
	return ChannelPairingService{Options: slices.Clone(opts)}
}

// Complete Channel pairing
func (r *ChannelPairingService) New(ctx context.Context, params ChannelPairingNewParams, opts ...option.RequestOption) (res *ChannelPairing, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "channel_pairings"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ChannelPairingNewParams struct {
	// Pairing code. Supply exactly one of Code or ChannelPairingID.
	Code string `json:"code,omitzero"`
	// Pending pairing ID. Supply exactly one of Code or ChannelPairingID.
	ChannelPairingID string `json:"channel_pairing_id,omitzero"`
	// Forward Identity ID to bind.
	IdentityID string `json:"identity_id" api:"required"`
	// Forward Template ID to bind.
	TemplateID string `json:"template_id" api:"required"`
	// Client-generated unique idempotency key, used to safely retry the same pairing
	// request.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ChannelPairingNewParams) MarshalJSON() ([]byte, error) {
	type shadow ChannelPairingNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChannelPairingNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Remove Channel pairing
func (r *ChannelPairingService) Delete(ctx context.Context, pairingID string, opts ...option.RequestOption) (res *DeletedChannelPairing, err error) {
	if pairingID == "" {
		return nil, fmt.Errorf("missing required pairing_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("channel_pairings/%s", url.PathEscape(pairingID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type DeletedChannelPairing struct {
	// Pairing ID.
	ID string `json:"id"`
	// Whether the pairing was removed; `true` on success.
	Deleted bool `json:"deleted"`
	JSON    struct {
		ID          respjson.Field
		Deleted     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DeletedChannelPairing) RawJSON() string { return r.JSON.raw }
func (r *DeletedChannelPairing) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChannelPairing struct {
	ScopeType        string `json:"scope_type"`
	ScopeExternalID  string `json:"scope_external_id"`
	ScopeDisplayName string `json:"scope_display_name"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	// Pairing ID, used when removing the pairing.
	ID string `json:"id"`
	// Always `channel_pairing`.
	Type string `json:"type"`
	// Channel ID.
	ChannelID string `json:"channel_id"`
	// Bound Forward Identity ID.
	IdentityID string `json:"identity_id"`
	// Bound Forward Template ID.
	TemplateID string `json:"template_id"`
	// `active` when the pairing succeeded.
	Status string `json:"status"`
	// Time the pairing completed.
	PairedAt string `json:"paired_at"`
	JSON     struct {
		ScopeType        respjson.Field
		ScopeExternalID  respjson.Field
		ScopeDisplayName respjson.Field
		CreatedAt        respjson.Field
		UpdatedAt        respjson.Field
		ID               respjson.Field
		Type             respjson.Field
		ChannelID        respjson.Field
		IdentityID       respjson.Field
		TemplateID       respjson.Field
		Status           respjson.Field
		PairedAt         respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r ChannelPairing) RawJSON() string                  { return r.JSON.raw }
func (r *ChannelPairing) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

func (r *ChannelPairingService) Get(ctx context.Context, pairingID string, opts ...option.RequestOption) (res *ChannelPairing, err error) {
	if pairingID == "" {
		return nil, fmt.Errorf("missing required pairing_id parameter")
	}
	opts = slices.Concat(r.Options, opts)
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, fmt.Sprintf("channel_pairings/%s", url.PathEscape(pairingID)), nil, &res, opts...)
	return res, err
}

// Update changes an active binding. Supply at least one nonempty field.
func (r *ChannelPairingService) Update(ctx context.Context, pairingID string, params ChannelPairingUpdateParams, opts ...option.RequestOption) (res *ChannelPairing, err error) {
	if pairingID == "" {
		return nil, fmt.Errorf("missing required pairing_id parameter")
	}
	opts = slices.Concat(r.Options, opts)
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, fmt.Sprintf("channel_pairings/%s", url.PathEscape(pairingID)), params, &res, opts...)
	return res, err
}
func (r *ChannelPairingService) List(ctx context.Context, params ChannelPairingListParams, opts ...option.RequestOption) (res *pagination.Page[ChannelPairing], err error) {
	opts = slices.Concat(r.Options, opts)
	var raw *http.Response
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	cfg, err := convention.NewRequestConfig(ctx, http.MethodGet, "channel_pairings", params, &res, opts...)
	if err != nil {
		return nil, err
	}
	if err = cfg.Execute(); err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}
func (r *ChannelPairingService) ListAutoPaging(ctx context.Context, params ChannelPairingListParams, opts ...option.RequestOption) *pagination.PageAutoPager[ChannelPairing] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type ChannelPairingUpdateParams struct {
	IdentityID param.Opt[string] `json:"identity_id,omitzero"`
	TemplateID param.Opt[string] `json:"template_id,omitzero"`
	paramObj
}

func (r ChannelPairingUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow ChannelPairingUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChannelPairingUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChannelPairingListParams struct {
	ChannelID        param.Opt[string] `query:"channel_id,omitzero" json:"-"`
	Status           param.Opt[string] `query:"status,omitzero" json:"-"`
	TargetType       param.Opt[string] `query:"target_type,omitzero" json:"-"`
	TargetExternalID param.Opt[string] `query:"target_external_id,omitzero" json:"-"`
	IdentityID       param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	TemplateID       param.Opt[string] `query:"template_id,omitzero" json:"-"`
	Limit            param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	AfterID          param.Opt[string] `query:"after_id,omitzero" json:"-"`
	BeforeID         param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r ChannelPairingListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{})
}
