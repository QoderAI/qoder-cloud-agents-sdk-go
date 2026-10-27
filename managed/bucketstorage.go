package managed

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// BucketStorageParam declares up to eight OSS bucket mounts in a cloud Environment.
type BucketStorageParam struct {
	Buckets []OSSBucketParam `json:"buckets" api:"required"`
	paramObj
}

func (r BucketStorageParam) MarshalJSON() ([]byte, error) {
	type shadow BucketStorageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BucketStorageParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// OSSBucketParam requires credentials for new buckets. Omit Credential on update to retain it.
type OSSBucketParam struct {
	Name                string                   `json:"name" api:"required"`
	Provider            string                   `json:"provider" api:"required"`
	Region              string                   `json:"region" api:"required"`
	Bucket              string                   `json:"bucket" api:"required"`
	Prefix              param.Opt[string]        `json:"prefix,omitzero"`
	Access              param.Opt[string]        `json:"access,omitzero"`
	AcceleratorEndpoint param.Opt[string]        `json:"accelerator_endpoint,omitzero"`
	Credential          OSSBucketCredentialParam `json:"credential,omitzero"`
	paramObj
}

func (r OSSBucketParam) MarshalJSON() ([]byte, error) {
	type shadow OSSBucketParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OSSBucketParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// OSSBucketCredentialParam accepts ram_role with RoleArn, or access_key with ID and secret.
// A read-view credential without secrets may be echoed to retain it on update.
type OSSBucketCredentialParam struct {
	Type            string            `json:"type" api:"required"`
	RoleArn         param.Opt[string] `json:"role_arn,omitzero"`
	AccessKeyID     param.Opt[string] `json:"access_key_id,omitzero"`
	AccessKeySecret param.Opt[string] `json:"access_key_secret,omitzero"`
	paramObj
}

func (r OSSBucketCredentialParam) MarshalJSON() ([]byte, error) {
	type shadow OSSBucketCredentialParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OSSBucketCredentialParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BucketStorage struct {
	Buckets []OSSBucket `json:"buckets"`
	JSON    struct {
		Buckets     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r BucketStorage) RawJSON() string                  { return r.JSON.raw }
func (r *BucketStorage) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type OSSBucket struct {
	Name                string              `json:"name"`
	Provider            string              `json:"provider"`
	Region              string              `json:"region"`
	Bucket              string              `json:"bucket"`
	Prefix              string              `json:"prefix"`
	Access              string              `json:"access"`
	AcceleratorEndpoint string              `json:"accelerator_endpoint"`
	Credential          OSSBucketCredential `json:"credential"`
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

func (r OSSBucket) RawJSON() string                  { return r.JSON.raw }
func (r *OSSBucket) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// OSSBucketCredential is the API read view. It contains no role ARN or secret.
type OSSBucketCredential struct {
	Type        string `json:"type"`
	AccessKeyID string `json:"access_key_id"`
	JSON        struct {
		Type        respjson.Field
		AccessKeyID respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r OSSBucketCredential) RawJSON() string                  { return r.JSON.raw }
func (r *OSSBucketCredential) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
