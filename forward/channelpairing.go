package forward

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
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
	// Six-digit pairing code shown in the Channel message.
	Code string `json:"code" api:"required"`
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
		ID          respjson.Field
		Type        respjson.Field
		ChannelID   respjson.Field
		IdentityID  respjson.Field
		TemplateID  respjson.Field
		Status      respjson.Field
		PairedAt    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ChannelPairing) RawJSON() string                  { return r.JSON.raw }
func (r *ChannelPairing) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
