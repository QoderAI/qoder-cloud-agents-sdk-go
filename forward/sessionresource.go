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

// 添加 Session 资源.
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
	// 资源类型，必须为 `file`。
	Type string `json:"type" api:"required"`
	// Files API 返回的 File ID，文件必须已上传完成。
	FileID string `json:"file_id" api:"required"`
	// Agent 容器内挂载路径；省略时由 Forward 根据文件名生成，默认挂载到 `/data/workspace/<文件名>`。
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
	// Session 资源 ID，以 `sesr_` 为前缀。
	ID string `json:"id"`
	// 资源类型，固定为 `file`。
	Type string `json:"type"`
	// 挂载的 File ID。
	FileID string `json:"file_id"`
	// 文件在 Agent 容器内的实际挂载路径。
	MountPath string `json:"mount_path"`
	// 资源创建时间，RFC3339。
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// 资源更新时间，RFC3339。
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
