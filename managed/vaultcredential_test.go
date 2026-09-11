package managed_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func TestVaultCredentialNew(t *testing.T) {
	client := contractClient(t, "VaultCredentialService", "New")
	response, err := client.Vaults.Credentials.New(context.Background(), "segment /?%#", managed.VaultCredentialNewParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "VaultCredential.New")
}

func TestVaultCredentialGet(t *testing.T) {
	client := contractClient(t, "VaultCredentialService", "Get")
	response, err := client.Vaults.Credentials.Get(context.Background(), "segment /?%#", managed.VaultCredentialGetParams{VaultID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "VaultCredential.Get")
}

func TestVaultCredentialUpdate(t *testing.T) {
	client := contractClient(t, "VaultCredentialService", "Update")
	response, err := client.Vaults.Credentials.Update(context.Background(), "segment /?%#", managed.VaultCredentialUpdateParams{VaultID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "VaultCredential.Update")
}

func TestVaultCredentialUpdateRejectsDisplayNameLocally(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value param.Opt[string]
	}{
		{"name", managed.String("unsupported name")},
		{"empty", managed.String("")},
		{"null", param.Null[string]()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				return reply(r, 200, `{}`), nil
			})
			params := managed.VaultCredentialUpdateParams{VaultID: "vault_test", DisplayName: tc.value}
			_, err := client.Vaults.Credentials.Update(context.Background(), "vcred_test", params)
			if err == nil || !strings.Contains(err.Error(), "display_name") || !strings.Contains(err.Error(), "not support") {
				t.Fatalf("expected unsupported display_name error, got %v", err)
			}
			if calls != 0 {
				t.Fatalf("unsupported update reached HTTP transport %d times", calls)
			}
		})
	}
}

func TestVaultCredentialUpdateSupportedBody(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != http.MethodPost || r.URL.EscapedPath() != "/api/v1/cloud/vaults/vault_test/credentials/vcred_test" {
			t.Fatalf("unexpected credential update route: %s %s", r.Method, r.URL.EscapedPath())
		}
		var got map[string]any
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		want := map[string]any{
			"auth":     map[string]any{"type": "static_bearer", "token": "rotated-test-token"},
			"metadata": map[string]any{"suite": "updated", "remove": nil},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("update body: got %#v want %#v", got, want)
		}
		return reply(r, 200, `{"id":"vcred_test","metadata":{"suite":"updated"}}`), nil
	})
	response, err := client.Vaults.Credentials.Update(context.Background(), "vcred_test", managed.VaultCredentialUpdateParams{
		VaultID: "vault_test",
		Auth: managed.VaultCredentialUpdateParamsAuthUnion{OfStaticBearer: &managed.ManagedAgentsStaticBearerUpdateParams{
			Type: "static_bearer", Token: managed.String("rotated-test-token"),
		}},
		Metadata: map[string]any{"suite": "updated", "remove": nil},
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 || response.ID != "vcred_test" || response.Metadata["suite"] != "updated" {
		t.Fatalf("credential update did not round trip: calls=%d", calls)
	}
}

func TestVaultCredentialList(t *testing.T) {
	client := contractClient(t, "VaultCredentialService", "List")
	response, err := client.Vaults.Credentials.List(context.Background(), "segment /?%#", managed.VaultCredentialListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "VaultCredential.List")
}

func TestVaultCredentialDelete(t *testing.T) {
	client := contractClient(t, "VaultCredentialService", "Delete")
	response, err := client.Vaults.Credentials.Delete(context.Background(), "segment /?%#", managed.VaultCredentialDeleteParams{VaultID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "VaultCredential.Delete")
}

func TestVaultCredentialArchive(t *testing.T) {
	client := contractClient(t, "VaultCredentialService", "Archive")
	response, err := client.Vaults.Credentials.Archive(context.Background(), "segment /?%#", managed.VaultCredentialArchiveParams{VaultID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "VaultCredential.Archive")
}

func TestVaultCredentialMCPOAuthValidate(t *testing.T) {
	client := contractClient(t, "VaultCredentialService", "MCPOAuthValidate")
	response, err := client.Vaults.Credentials.MCPOAuthValidate(context.Background(), "segment /?%#", managed.VaultCredentialMCPOAuthValidateParams{VaultID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "VaultCredential.MCPOAuthValidate")
}

func TestTypedUserEventAndCredentialUpdate(t *testing.T) {
	var bodies []map[string]any
	c := testClient(func(r *http.Request) (*http.Response, error) {
		var m map[string]any
		if e := json.NewDecoder(r.Body).Decode(&m); e != nil {
			t.Fatal(e)
		}
		bodies = append(bodies, m)
		return reply(r, 200, `{}`), nil
	})
	var params managed.SessionEventSendParams
	if e := json.Unmarshal([]byte(`{"events":[{"type":"user.message","content":[{"type":"text","text":"你好"}]}]}`), &params); e != nil {
		t.Fatal(e)
	}
	if _, e := c.Sessions.Events.Send(context.Background(), "session", params); e != nil {
		t.Fatal(e)
	}
	if bodies[0]["events"].([]any)[0].(map[string]any)["content"].([]any)[0].(map[string]any)["text"] != "你好" {
		t.Fatal("typed event union lost content")
	}
	var p managed.VaultCredentialUpdateParams
	if e := json.Unmarshal([]byte(`{"auth":{"type":"mcp_oauth","access_token":"rotated","expires_at":null},"metadata":{"team":"sdk"}}`), &p); e != nil {
		t.Fatal(e)
	}
	p.VaultID = "vault"
	if _, e := c.Vaults.Credentials.Update(context.Background(), "cred", p); e != nil {
		t.Fatal(e)
	}
	auth := bodies[1]["auth"].(map[string]any)
	if auth["type"] != "mcp_oauth" || auth["access_token"] != "rotated" {
		t.Fatal(auth)
	}
	if v, ok := auth["expires_at"]; !ok || v != nil {
		t.Fatal("credential expiry clear lost")
	}
	if _, e := c.Vaults.Credentials.MCPOAuthValidate(context.Background(), "cred", managed.VaultCredentialMCPOAuthValidateParams{VaultID: "vault"}); e != nil {
		t.Fatal(e)
	}
	if len(bodies[2]) != 0 {
		t.Fatal("validation requires empty JSON object")
	}
}
