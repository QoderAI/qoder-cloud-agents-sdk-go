package forward

import (
	"encoding/json"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

type EnvironmentConfigParams struct {
	Type          param.Opt[string]              `json:"type,omitzero"`
	Packages      EnvironmentConfigPackagesParam `json:"packages,omitzero"`
	SetupScript   param.Opt[string]              `json:"setup_script,omitzero"`
	StartupScript param.Opt[string]              `json:"startup_script,omitzero"`
	Storage       EnvironmentStorageParam        `json:"storage,omitzero"`
	configParamObj
}

func (r EnvironmentConfigParams) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentConfigParams
	return marshalConfigParam(r, (*shadow)(&r))
}
func (r *EnvironmentConfigParams) UnmarshalAPIJSON(data []byte) error { return r.UnmarshalJSON(data) }
func (r *EnvironmentConfigParams) UnmarshalJSON(data []byte) error {
	return unmarshalObjectParam[EnvironmentConfigParams](data, r)
}

func unmarshalEnvironmentConfig(data []byte, dst *EnvironmentConfigParams) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if raw, present := fields["config"]; present {
		return dst.UnmarshalJSON(raw)
	}
	return nil
}

type EnvironmentConfigPackagesParam struct {
	Apt   []string `json:"apt,omitzero"`
	Npm   []string `json:"npm,omitzero"`
	Pip   []string `json:"pip,omitzero"`
	Cargo []string `json:"cargo,omitzero"`
	Gem   []string `json:"gem,omitzero"`
	Go    []string `json:"go,omitzero"`
	configParamObj
}

func (r EnvironmentConfigPackagesParam) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentConfigPackagesParam
	return marshalConfigParam(r, (*shadow)(&r))
}
func (r *EnvironmentConfigPackagesParam) UnmarshalAPIJSON(data []byte) error {
	return r.UnmarshalJSON(data)
}
func (r *EnvironmentConfigPackagesParam) UnmarshalJSON(data []byte) error {
	return unmarshalObjectParam[EnvironmentConfigPackagesParam](data, r)
}

type EnvironmentStorageParam struct {
	Buckets []EnvironmentStorageBucketParam `json:"buckets,omitzero"`
	configParamObj
}

func (r EnvironmentStorageParam) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentStorageParam
	return marshalConfigParam(r, (*shadow)(&r))
}
func (r *EnvironmentStorageParam) UnmarshalAPIJSON(data []byte) error { return r.UnmarshalJSON(data) }
func (r *EnvironmentStorageParam) UnmarshalJSON(data []byte) error {
	return unmarshalObjectParam[EnvironmentStorageParam](data, r)
}

type EnvironmentStorageBucketParam struct {
	Name                string                            `json:"name" api:"required"`
	Provider            string                            `json:"provider" api:"required"`
	Region              string                            `json:"region" api:"required"`
	Bucket              string                            `json:"bucket" api:"required"`
	Prefix              param.Opt[string]                 `json:"prefix,omitzero"`
	Access              param.Opt[string]                 `json:"access,omitzero"`
	AcceleratorEndpoint param.Opt[string]                 `json:"accelerator_endpoint,omitzero"`
	Credential          EnvironmentStorageCredentialParam `json:"credential,omitzero"`
	configParamObj
}

func (r EnvironmentStorageBucketParam) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentStorageBucketParam
	return marshalConfigParam(r, (*shadow)(&r))
}
func (r *EnvironmentStorageBucketParam) UnmarshalAPIJSON(data []byte) error {
	return r.UnmarshalJSON(data)
}
func (r *EnvironmentStorageBucketParam) UnmarshalJSON(data []byte) error {
	return unmarshalObjectParam[EnvironmentStorageBucketParam](data, r)
}

type EnvironmentStorageCredentialParam struct {
	Type            string            `json:"type" api:"required"`
	RoleARN         param.Opt[string] `json:"role_arn,omitzero"`
	AccessKeyID     param.Opt[string] `json:"access_key_id,omitzero"`
	AccessKeySecret param.Opt[string] `json:"access_key_secret,omitzero"`
	configParamObj
}

func (r EnvironmentStorageCredentialParam) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentStorageCredentialParam
	return marshalConfigParam(r, (*shadow)(&r))
}
func (r *EnvironmentStorageCredentialParam) UnmarshalAPIJSON(data []byte) error {
	return r.UnmarshalJSON(data)
}
func (r *EnvironmentStorageCredentialParam) UnmarshalJSON(data []byte) error {
	return unmarshalObjectParam[EnvironmentStorageCredentialParam](data, r)
}

type EnvironmentStorage struct {
	Buckets []EnvironmentStorageBucket `json:"buckets"`
	JSON    struct {
		Buckets     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EnvironmentStorage) RawJSON() string                  { return r.JSON.raw }
func (r *EnvironmentStorage) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type EnvironmentStorageBucket struct {
	Name                string                       `json:"name"`
	Provider            string                       `json:"provider"`
	Region              string                       `json:"region"`
	Bucket              string                       `json:"bucket"`
	Prefix              string                       `json:"prefix"`
	Access              string                       `json:"access"`
	AcceleratorEndpoint string                       `json:"accelerator_endpoint"`
	Credential          EnvironmentStorageCredential `json:"credential"`
	JSON                struct {
		Name                respjson.Field
		Provider            respjson.Field
		Region              respjson.Field
		Bucket              respjson.Field
		Prefix              respjson.Field
		Access              respjson.Field
		AcceleratorEndpoint respjson.Field
		Credential          respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

func (r EnvironmentStorageBucket) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentStorageBucket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EnvironmentStorageCredential struct {
	Type        string `json:"type"`
	AccessKeyID string `json:"access_key_id"`
	JSON        struct {
		Type        respjson.Field
		AccessKeyID respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EnvironmentStorageCredential) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentStorageCredential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
