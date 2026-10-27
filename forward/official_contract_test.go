package forward_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestOfficialOutcomeRubricContract(t *testing.T) {
	for _, rubric := range []forward.OutcomeRubricParam{{Type: "text", Content: forward.String("report exists")}, {Type: "file", FileID: forward.String("file_test")}} {
		wire, err := json.Marshal(forward.SessionEventParam{Type: "user.define_outcome", Rubric: rubric})
		if err != nil {
			t.Fatal(err)
		}
		var event forward.SessionEvent
		if err = json.Unmarshal(wire, &event); err != nil {
			t.Fatal(err)
		}
		if !event.JSON.Rubric.Valid() || event.Rubric.Type != rubric.Type {
			t.Fatalf("invalid rubric: %+v", event)
		}
		if rubric.Type == "text" && event.Rubric.Content != "report exists" {
			t.Fatal(event.Rubric)
		}
		if rubric.Type == "file" && event.Rubric.FileID != "file_test" {
			t.Fatal(event.Rubric)
		}
	}
}

func TestOfficialNewAPIContracts(t *testing.T) {
	for _, tc := range []struct {
		id   string
		call func(forward.Client) (any, error)
	}{
		{"createDream", func(c forward.Client) (any, error) {
			return c.Dreams.New(context.Background(), contractParams[forward.DreamNewParams](t, "createDream"))
		}},
		{"listDreams", func(c forward.Client) (any, error) {
			return c.Dreams.List(context.Background(), contractParams[forward.DreamListParams](t, "listDreams"))
		}},
		{"getDream", func(c forward.Client) (any, error) { return c.Dreams.Get(context.Background(), pathSegment) }},
		{"cancelDream", func(c forward.Client) (any, error) { return c.Dreams.Cancel(context.Background(), pathSegment) }},
		{"archiveDream", func(c forward.Client) (any, error) { return c.Dreams.Archive(context.Background(), pathSegment) }},
		{"listChannelPairings", func(c forward.Client) (any, error) {
			return c.ChannelPairings.List(context.Background(), contractParams[forward.ChannelPairingListParams](t, "listChannelPairings"))
		}},
		{"getChannelPairing", func(c forward.Client) (any, error) { return c.ChannelPairings.Get(context.Background(), pathSegment) }},
		{"updateChannelPairing", func(c forward.Client) (any, error) {
			return c.ChannelPairings.Update(context.Background(), pathSegment, contractParams[forward.ChannelPairingUpdateParams](t, "updateChannelPairing"))
		}},
	} {
		t.Run(tc.id, func(t *testing.T) {
			res, err := tc.call(contractClient(t, tc.id, 200))
			if err != nil {
				t.Fatal(err)
			}
			checkDecoded(t, res)
		})
	}
}

func TestOfficialIdentityNullAndArrayOverrides(t *testing.T) {
	for _, tc := range []struct {
		spec forward.IdentityConfigSpecParam
		want string
	}{
		{forward.IdentityConfigSpecParam{Tools: forward.IdentityToolsUnionParam{OfArray: []forward.ToolParam{}}}, `{"tools":[]}`},
		{forward.IdentityConfigSpecParam{Tools: forward.IdentityToolsUnionParam{OfOverrides: map[string]forward.ToolOverrideParam{"Bash": {Enabled: forward.Bool(false)}}}}, `{"tools":{"Bash":{"enabled":false}}}`},
		{forward.IdentityConfigSpecParam{Tools: param.NullStruct[forward.IdentityToolsUnionParam](), ManagedToolConfig: param.NullMap[map[string]forward.ManagedCapabilityOverrideParam]()}, `{"tools":null,"managed_tool_config":null}`},
		{forward.IdentityConfigSpecParam{ManagedToolConfig: map[string]forward.ManagedCapabilityOverrideParam{"schedule": {Enabled: false}, "drive": param.NullStruct[forward.ManagedCapabilityOverrideParam]()}}, `{"managed_tool_config":{"schedule":{"enabled":false},"drive":null}}`},
	} {
		var want any
		_ = json.Unmarshal([]byte(tc.want), &want)
		body, err := json.Marshal(tc.spec)
		if err != nil {
			t.Fatal(err)
		}
		var got any
		_ = json.Unmarshal(body, &got)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %s want %s", body, tc.want)
		}
		var spec forward.IdentityConfigSpec
		if err = json.Unmarshal(body, &spec); err != nil {
			t.Fatal(err)
		}
	}
	var roster forward.MultiagentConfig
	raw := `{"type":"roster","agents":["tmpl_helper",{"type":"advisor","model":"ultimate"},{"type":"self"}]}`
	if err := json.Unmarshal([]byte(raw), &roster); err != nil {
		t.Fatal(err)
	}
	if roster.Agents[0].OfString != "tmpl_helper" || roster.Agents[1].Model != "ultimate" {
		t.Fatal(roster)
	}
	var params forward.MultiagentConfigParam
	if err := json.Unmarshal([]byte(raw), &params); err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}
	var got, want any
	_ = json.Unmarshal(body, &got)
	_ = json.Unmarshal([]byte(raw), &want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("roster: %s", body)
	}
}

func TestOfficialIdentityToolsJSONRoundTrip(t *testing.T) {
	for _, tools := range []string{
		`[]`,
		`[{"type":"agent_toolset_20260401"}]`,
		`{}`,
		`{"Bash":{"enabled":false}}`,
		`null`,
	} {
		t.Run(tools, func(t *testing.T) {
			var params forward.IdentityConfigUpsertParams
			input := `{"identity_config":{"tools":` + tools + `}}`
			if err := json.Unmarshal([]byte(input), &params); err != nil {
				t.Fatal(err)
			}
			client := testClient(func(req *http.Request) (*http.Response, error) {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}
				var got, want any
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(input), &want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("request changed from %s to %s", input, body)
				}
				return reply(req, 200, `{}`), nil
			})
			if _, err := client.Identities.Configs.Upsert(context.Background(), "idn", "tmpl", params); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestOfficialMountPaginationRetainsScope(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		calls := 0
		key := "after_id"
		if reverse {
			key = "before_id"
		}
		client := testClient(func(r *http.Request) (*http.Response, error) {
			calls++
			q := r.URL.Query()
			if r.URL.Path != "/api/v1/forward/identities/idn/templates/tmpl/memory_stores" || q.Get("limit") != "1" {
				t.Fatal(r.URL)
			}
			want := "start"
			if calls == 2 {
				want = "mem_one"
			}
			if q.Get(key) != want {
				t.Fatal(r.URL)
			}
			id := "mem_one"
			more := true
			if calls == 2 {
				id = "mem_two"
				more = false
			}
			return reply(r, 200, fmt.Sprintf(`{"data":[{"memory_store_id":%q,"access":"read_write"}],"first_id":%q,"last_id":%q,"has_more":%t}`, id, id, id, more)), nil
		})
		p := forward.IdentityMemoryStoreListParams{Limit: forward.Int(1)}
		if reverse {
			p.BeforeID = forward.String("start")
		} else {
			p.AfterID = forward.String("start")
		}
		pager := client.Identities.MemoryStores.ListAutoPaging(context.Background(), "idn", "tmpl", p)
		var ids []string
		for pager.Next() {
			ids = append(ids, pager.Current().MemoryStoreID)
		}
		if pager.Err() != nil || !reflect.DeepEqual(ids, []string{"mem_one", "mem_two"}) || calls != 2 {
			t.Fatalf("%v %v calls=%d", ids, pager.Err(), calls)
		}
	}
}

func TestOfficialForwardParameterPlacement(t *testing.T) {
	client := testClient(func(r *http.Request) (*http.Response, error) {
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
		}
		switch r.URL.Path {
		case "/api/v1/forward/identities/idn":
			if string(body) != `{"auto_memory_store":false}` {
				t.Fatal(string(body))
			}
			return reply(r, 200, `{"id":"idn","auto_memory_store":false}`), nil
		case "/api/v1/forward/identities/idn/templates/tmpl/memory_stores":
			var b map[string]any
			_ = json.Unmarshal(body, &b)
			if b["access"] != "read_write" || b["memory_store_id"] != "mem" {
				t.Fatal(b)
			}
			return reply(r, 200, `{"memory_store_id":"mem","access":"read_write"}`), nil
		case "/api/v1/forward/schedules/sched":
			var b map[string]any
			_ = json.Unmarshal(body, &b)
			if len(b) != 1 || b["resources"] == nil {
				t.Fatal(b)
			}
			return reply(r, 200, `{"resources":[{"type":"file","file_id":"file_test","mount_path":"/data/test.txt"}]}`), nil
		case "/api/v1/forward/schedule_runs/run":
			if r.URL.Query().Get("include_initial_events") != "true" || len(body) != 0 {
				t.Fatalf("%s %s", r.URL, body)
			}
			return reply(r, 200, `{"initial_events":[{"type":"user.message","content":[{"type":"text","text":"run input"}]}]}`), nil
		}
		t.Fatal(r.URL)
		return nil, nil
	})
	ctx := context.Background()
	id, err := client.Identities.Update(ctx, "idn", forward.IdentityUpdateParams{AutoMemoryStore: forward.Bool(false)})
	if err != nil || !id.JSON.AutoMemoryStore.Valid() || id.AutoMemoryStore {
		t.Fatalf("%+v %v", id, err)
	}
	if _, err = client.Identities.MemoryStores.Mount(ctx, "idn", "tmpl", forward.IdentityMemoryStoreMountParams{MemoryStoreID: "mem", Access: forward.String("read_write")}); err != nil {
		t.Fatal(err)
	}
	sched, err := client.Schedules.Update(ctx, "sched", forward.ScheduleUpdateParams{Resources: []forward.ScheduleFileResourceParam{{Type: "file", FileID: "file_test", MountPath: forward.String("/data/test.txt")}}})
	if err != nil || len(sched.Resources) != 1 || !sched.JSON.Resources.Valid() {
		t.Fatalf("%+v %v", sched, err)
	}
	run, err := client.ScheduleRuns.Get(ctx, "run", forward.ScheduleRunGetParams{IncludeInitialEvents: forward.Bool(true)})
	if err != nil || len(run.InitialEvents) != 1 || !run.JSON.InitialEvents.Valid() {
		t.Fatalf("%+v %v", run, err)
	}
}
