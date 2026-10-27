package managed_test

import (
	"context"
	"encoding/json"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestOfficialDreamModelContract(t *testing.T) {
	body := jsonObject(t, managed.DreamNewParams{Inputs: []managed.DreamInputUnionParam{}})
	if _, ok := body["model"]; ok {
		t.Fatal("default model should be omitted", body)
	}
	cfg := managed.DreamModelConfigParam{ID: "auto", Effort: managed.String("none"), ContextWindow: managed.Int(400000)}
	body = jsonObject(t, cfg)
	if body["effort"] != "none" || body["context_window"] != float64(400000) {
		t.Fatal(body)
	}
	for _, raw := range []string{`{"model":{"id":"auto","effort":"none","context_window":400000}}`} {
		var dream managed.Dream
		if err := json.Unmarshal([]byte(raw), &dream); err != nil {
			t.Fatal(err)
		}
		if dream.Model.ID != "auto" || dream.Model.Effort != "none" || dream.Model.ContextWindow != 400000 {
			t.Fatalf("lost model: %+v", dream.Model)
		}
	}
}

func TestOfficialOAuthAPIContract(t *testing.T) {
	c := contractClient(t, "OAuthService", "Start")
	result, err := c.OAuth.Start(context.Background(), managed.OAuthStartParams{VaultID: "vault", MCPServerURL: "https://mcp.example.com"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(*result), "oauth")
}
func TestOfficialOAuthSecretAndMemoryCursors(t *testing.T) {
	calls := 0
	c := managed.NewClient(option.WithPAT("test"), option.WithBaseURL("https://qoder.test/api/v1/cloud"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			if r.URL.Path != "/api/v1/cloud/oauth/start" || r.Method != "POST" || r.URL.RawQuery != "" || r.Header.Get("qoder-workspace-id") != "workspace" {
				t.Fatal(r.URL)
			}
			b, _ := io.ReadAll(r.Body)
			var body map[string]any
			_ = json.Unmarshal(b, &body)
			if !reflect.DeepEqual(body, map[string]any{"vault_id": "vault", "mcp_server_url": "https://mcp.example.com", "client_id": "client", "client_secret": "test-secret"}) {
				t.Fatal(body)
			}
			return reply(r, 200, `{"authorization_url":"https://auth.example.com","state":"opaque","callback_origin":"https://api.example.com"}`), nil
		}
		if r.URL.Path != "/api/v1/cloud/memory_stores" || r.URL.Query().Get("before_id") != "mem_before" || r.URL.Query().Get("include_archived") != "false" {
			t.Fatal(r.URL)
		}
		return reply(r, 200, `{"data":[],"has_more":false}`), nil
	})}))
	result, err := c.OAuth.Start(context.Background(), managed.OAuthStartParams{VaultID: "vault", MCPServerURL: "https://mcp.example.com", ClientID: managed.String("client"), ClientSecret: managed.String("test-secret"), WorkspaceID: managed.String("workspace")})
	if err != nil || result.State != "opaque" || !result.JSON.AuthorizationURL.Valid() {
		t.Fatalf("%+v %v", result, err)
	}
	if _, err = c.MemoryStores.List(context.Background(), managed.MemoryStoreListParams{BeforeID: managed.String("mem_before"), IncludeArchived: managed.Bool(false)}); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatal(calls)
	}
}
func TestOfficialOSSStorageWriteReadAndClear(t *testing.T) {
	for _, bucket := range []managed.OSSBucketParam{
		{Name: "assets", Provider: "aliyun-oss", Region: "cn-hangzhou", Bucket: "example", Credential: managed.OSSBucketCredentialParam{Type: "ram_role", RoleArn: managed.String("test-role")}},
		{Name: "assets", Provider: "aliyun-oss", Region: "cn-hangzhou", Bucket: "example", Credential: managed.OSSBucketCredentialParam{Type: "access_key", AccessKeyID: managed.String("test-ak"), AccessKeySecret: managed.String("test-secret")}},
		{Name: "assets", Provider: "aliyun-oss", Region: "cn-hangzhou", Bucket: "example"},
	} {
		body := jsonObject(t, managed.CloudConfigParams{Storage: managed.BucketStorageParam{Buckets: []managed.OSSBucketParam{bucket}}})
		stored := body["storage"].(map[string]any)["buckets"].([]any)[0].(map[string]any)
		if bucket.Credential.Type == "" {
			if _, ok := stored["credential"]; ok {
				t.Fatal("must preserve stored credentials by omission", stored)
			}
		} else {
			cred := stored["credential"].(map[string]any)
			if cred["type"] != bucket.Credential.Type {
				t.Fatal(cred)
			}
			if bucket.Credential.Type == "access_key" && cred["access_key_secret"] != "test-secret" {
				t.Fatal(cred)
			}
		}
	}
	body := jsonObject(t, managed.CloudConfigParams{Storage: managed.BucketStorageParam{Buckets: []managed.OSSBucketParam{}}})
	if len(body["storage"].(map[string]any)["buckets"].([]any)) != 0 {
		t.Fatal(body)
	}
	var env managed.Environment
	if err := json.Unmarshal([]byte(`{"config":{"type":"cloud","storage":{"buckets":[{"name":"assets","provider":"aliyun-oss","region":"cn-hangzhou","bucket":"example","credential":{"type":"access_key","access_key_id":"test****"}}]}}}`), &env); err != nil {
		t.Fatal(err)
	}
	if !env.Config.JSON.Storage.Valid() || env.Config.Storage.Buckets[0].Credential.AccessKeyID != "test****" {
		t.Fatalf("%+v", env.Config)
	}
	if _, ok := reflect.TypeOf(managed.OSSBucketCredential{}).FieldByName("AccessKeySecret"); ok {
		t.Fatal("read view must not expose secret")
	}
}
