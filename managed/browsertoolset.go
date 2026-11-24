package managed

import (
	"encoding/json"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/constant"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

type ManagedBrowserToolset20260714Type = constant.BrowserToolset20260714

const ManagedBrowserToolset20260714TypeBrowserToolset20260714 ManagedBrowserToolset20260714Type = "browser_toolset_20260714"

type ManagedBrowserToolset20260714Param struct {
	Type ManagedBrowserToolset20260714Type `json:"type" default:"browser_toolset_20260714" api:"required"`
	paramObj
}

func (r ManagedBrowserToolset20260714Param) MarshalJSON() ([]byte, error) {
	type shadow ManagedBrowserToolset20260714Param
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedBrowserToolset20260714Param) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedBrowserToolset20260714 struct {
	Type ManagedBrowserToolset20260714Type `json:"type" api:"required"`
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ManagedBrowserToolset20260714) RawJSON() string { return r.JSON.raw }
func (r *ManagedBrowserToolset20260714) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r ManagedBrowserToolset20260714) ToParam() ManagedBrowserToolset20260714Param {
	return param.Override[ManagedBrowserToolset20260714Param](json.RawMessage(r.RawJSON()))
}

func (ManagedBrowserToolset20260714) implManagedAgentsAgentToolUnion()              {}
func (ManagedBrowserToolset20260714) implManagedAgentsSessionAgentToolUnion()       {}
func (ManagedBrowserToolset20260714) implManagedAgentsSessionThreadAgentToolUnion() {}

func (u ManagedAgentsAgentToolUnion) AsBrowserToolset20260714() (v ManagedBrowserToolset20260714) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsSessionAgentToolUnion) AsBrowserToolset20260714() (v ManagedBrowserToolset20260714) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsSessionThreadAgentToolUnion) AsBrowserToolset20260714() (v ManagedBrowserToolset20260714) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
