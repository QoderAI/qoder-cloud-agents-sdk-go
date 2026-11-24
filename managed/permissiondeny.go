package managed

import (
	"encoding/json"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/constant"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

type ManagedAgentsAlwaysDenyPolicyType = constant.AlwaysDeny

const ManagedAgentsAlwaysDenyPolicyTypeAlwaysDeny ManagedAgentsAlwaysDenyPolicyType = "always_deny"

type ManagedAgentsAlwaysDenyPolicyParam struct {
	Type ManagedAgentsAlwaysDenyPolicyType `json:"type" default:"always_deny" api:"required"`
	paramObj
}

func (r ManagedAgentsAlwaysDenyPolicyParam) MarshalJSON() ([]byte, error) {
	type shadow ManagedAgentsAlwaysDenyPolicyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsAlwaysDenyPolicyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAlwaysDenyPolicy struct {
	Type ManagedAgentsAlwaysDenyPolicyType `json:"type" api:"required"`
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ManagedAgentsAlwaysDenyPolicy) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAlwaysDenyPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r ManagedAgentsAlwaysDenyPolicy) ToParam() ManagedAgentsAlwaysDenyPolicyParam {
	return param.Override[ManagedAgentsAlwaysDenyPolicyParam](json.RawMessage(r.RawJSON()))
}

func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsBashToolConfigPermissionPolicyUnion()      {}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsEditToolConfigPermissionPolicyUnion()      {}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsReadToolConfigPermissionPolicyUnion()      {}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsWriteToolConfigPermissionPolicyUnion()     {}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsGlobToolConfigPermissionPolicyUnion()      {}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsGrepToolConfigPermissionPolicyUnion()      {}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsWebFetchToolConfigPermissionPolicyUnion()  {}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsWebSearchToolConfigPermissionPolicyUnion() {}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsMCPToolConfigPermissionPolicyUnion()       {}
func (ManagedAgentsAlwaysDenyPolicy) implManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion() {
}

func (u ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsBashToolConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsEditToolConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsReadToolConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsWriteToolConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsGlobToolConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsGrepToolConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsWebFetchToolConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsWebSearchToolConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsMCPToolConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysDeny() (v ManagedAgentsAlwaysDenyPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
