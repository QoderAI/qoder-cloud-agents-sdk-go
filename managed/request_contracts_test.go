package managed_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
)

func TestManagedDocumentedRequestBodies(t *testing.T) {
	var cases []struct {
		Service, Method, Doc string
		Case                 int
		Body                 json.RawMessage
	}
	b, err := os.ReadFile("testdata/request-bodies.json")
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) != 23 {
		t.Fatalf("documented body cases=%d", len(cases))
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s/%s/%d", c.Service, c.Method, c.Case), func(t *testing.T) {
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Body == nil {
					t.Fatal("missing request body")
				}
				actual, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				var want, got any
				if err = json.Unmarshal(c.Body, &want); err != nil {
					t.Fatal(err)
				}
				if err = json.Unmarshal(actual, &got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("%s request body\n got: %s\nwant: %s", c.Doc, actual, c.Body)
				}
				return reply(r, 200, `{}`), nil
			})
			_, err := testutil.InvokeJSON(t, client, testutil.Endpoint{Service: c.Service, Method: c.Method}, c.Body)
			if err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Fatalf("calls=%d", calls)
			}
		})
	}
}

func TestManagedNullableMetadataAndRuntimeParameters(t *testing.T) {
	cases := []struct {
		name string
		body any
		want string
	}{
		{"nullable_metadata", managed.SessionUpdateParams{Metadata: map[string]any{"remove": nil, "empty": "", "keep": "value"}}, `{"metadata":{"remove":null,"empty":"","keep":"value"}}`},
		{"tool_allow_and_deny", managed.ManagedAgentsAgentToolset20260401Params{Type: "agent_toolset_20260401", EnabledTools: []string{"Read", "Bash"}, DisallowedTools: []string{"WebFetch"}}, `{"type":"agent_toolset_20260401","enabled_tools":["Read","Bash"],"disallowed_tools":["WebFetch"]}`},
		{"setup_script", managed.CloudConfigParams{SetupScript: managed.String("printf ready")}, `{"type":"cloud","setup_script":"printf ready"}`},
		{"clear_setup_script", managed.CloudConfigParams{SetupScript: param.Null[string]()}, `{"type":"cloud","setup_script":null}`},
		{"file_mount", managed.SessionResourceAddParams{ManagedAgentsFileResourceParams: managed.ManagedAgentsFileResourceParams{Type: "file", FileID: "file-one", MountPath: managed.String("/data/a b.txt")}}, `{"type":"file","file_id":"file-one","mount_path":"/data/a b.txt"}`},
		{"force_stop", managed.EnvironmentWorkStopParams{EnvironmentID: "env", SelfHostedWorkStopRequest: managed.SelfHostedWorkStopRequestParam{Force: managed.Bool(false)}}, `{"force":false}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := json.Marshal(c.body)
			if err != nil {
				t.Fatal(err)
			}
			var got, want any
			if err = json.Unmarshal(b, &got); err != nil {
				t.Fatal(err)
			}
			if err = json.Unmarshal([]byte(c.want), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %s want %s", b, c.want)
			}
		})
	}
}

func TestManagedQueryAndHeaderContract(t *testing.T) {
	when := time.Date(2026, 9, 9, 3, 4, 5, 0, time.UTC)
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		want := url.Values{"agent_id": {"agent /?"}, "agent_version": {"0"}, "include_archived": {"false"}, "limit": {"100"}, "page": {"page +/="}, "created_at[gte]": {"2026-09-09T03:04:05Z"}, "statuses[]": {"idle", "running"}, "order": {"asc"}}
		if !reflect.DeepEqual(r.URL.Query(), want) {
			t.Errorf("query=%v want=%v", r.URL.Query(), want)
		}
		if !reflect.DeepEqual(r.Header.Values("X-Qoder-Beta"), []string{"first", "second"}) {
			t.Error("typed headers missing")
		}
		return reply(r, 200, `{"data":[],"next_page":null}`), nil
	})
	_, err := client.Sessions.List(context.Background(), managed.SessionListParams{AgentID: managed.String("agent /?"), AgentVersion: managed.Int(0), IncludeArchived: managed.Bool(false), Limit: managed.Int(100), Page: managed.String("page +/="), CreatedAtGte: param.NewOpt(when), Statuses: []string{"idle", "running"}, Order: "asc", Betas: []managed.QoderBeta{"first", "second"}})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("calls=%d", calls)
	}
}

func TestManagedRuntimeResponseFields(t *testing.T) {
	var agent managed.ManagedAgentsAgent
	if err := json.Unmarshal([]byte(`{"id":"agent","tools":[{"type":"agent_toolset_20260401","enabled_tools":["Read"],"disallowed_tools":["WebFetch"]}]}`), &agent); err != nil {
		t.Fatal(err)
	}
	if len(agent.Tools) != 1 || !reflect.DeepEqual(agent.Tools[0].EnabledTools, []string{"Read"}) || !reflect.DeepEqual(agent.Tools[0].DisallowedTools, []string{"WebFetch"}) {
		t.Fatal("tool allow/deny lists were not decoded")
	}
	var config managed.CloudConfig
	if err := json.Unmarshal([]byte(`{"type":"cloud","setup_script":"echo ready"}`), &config); err != nil {
		t.Fatal(err)
	}
	if config.SetupScript != "echo ready" || !config.JSON.SetupScript.Valid() {
		t.Fatal("setup_script response lost")
	}
}
