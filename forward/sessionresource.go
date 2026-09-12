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
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// SessionResourceService provides Forward SessionResource operations.
type SessionResourceService struct {
	Options []option.RequestOption
}

func NewSessionResourceService(opts ...option.RequestOption) SessionResourceService {
	return SessionResourceService{Options: slices.Clone(opts)}
}

// Add Session resource
func (r *SessionResourceService) Add(ctx context.Context, sessionID string, params SessionResourceAddParams, opts ...option.RequestOption) (res *SessionResource, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/resources", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type SessionResourceAddParams struct {
	// Resource type; must be `file`.
	Type string `json:"type" api:"required"`
	// File ID returned by the Files API; the file must already be fully uploaded.
	FileID string `json:"file_id" api:"required"`
	// Mount path inside the Agent container; when omitted, Forward derives one from the
	// file name and mounts it at `/data/workspace/<file name>`.
	MountPath param.Opt[string] `json:"mount_path,omitzero"`
	paramObj
}

func (r SessionResourceAddParams) MarshalJSON() ([]byte, error) {
	type shadow SessionResourceAddParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionResourceAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionResource struct {
	// Session resource ID, prefixed with `sesr_`.
	ID string `json:"id"`
	// Resource type, always `file`.
	Type string `json:"type"`
	// ID of the mounted File.
	FileID string `json:"file_id"`
	// Actual mount path of the file inside the Agent container.
	MountPath string `json:"mount_path"`
	// Resource creation time, RFC3339.
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// Resource update time, RFC3339.
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	JSON      struct {
		ID          respjson.Field
		Type        respjson.Field
		FileID      respjson.Field
		MountPath   respjson.Field
		CreatedAt   respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SessionResource) RawJSON() string                  { return r.JSON.raw }
func (r *SessionResource) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
