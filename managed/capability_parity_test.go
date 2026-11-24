package managed_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/constant"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func parityJSON(t *testing.T, value any, want string) {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var actual, expected any
	if err := json.Unmarshal(raw, &actual); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(want), &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got %s want %s", raw, want)
	}
}

func TestCapabilityParityDenyPolicies(t *testing.T) {
	deny := &managed.ManagedAgentsAlwaysDenyPolicyParam{Type: "always_deny"}
	for _, request := range []any{
		managed.ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsBashToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsEditToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsReadToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
		managed.ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny},
	} {
		t.Run(reflect.TypeOf(request).Name(), func(t *testing.T) {
			parityJSON(t, request, `{"type":"always_deny"}`)
			decoded := reflect.New(reflect.TypeOf(request))
			if err := json.Unmarshal([]byte(`{"type":"always_deny"}`), decoded.Interface()); err != nil {
				t.Fatal(err)
			}
			if decoded.Elem().FieldByName("OfAlwaysDeny").IsNil() {
				t.Fatal("missing deny discriminator")
			}
			parityJSON(t, decoded.Interface(), `{"type":"always_deny"}`)
			getType := decoded.MethodByName("GetType").Call(nil)[0].Interface().(*string)
			if getType == nil || *getType != "always_deny" {
				t.Fatal(getType)
			}
		})
	}
	for _, response := range []any{
		&managed.ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsBashToolConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsEditToolConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsReadToolConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsWriteToolConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsGlobToolConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsGrepToolConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsWebFetchToolConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsWebSearchToolConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsMCPToolConfigPermissionPolicyUnion{},
		&managed.ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion{},
	} {
		t.Run(reflect.TypeOf(response).Elem().Name(), func(t *testing.T) {
			for _, raw := range []string{`{"type":"always_deny","future":7}`, `{"type":"future_policy","future":7}`} {
				if err := json.Unmarshal([]byte(raw), response); err != nil {
					t.Fatal(err)
				}
				v := reflect.ValueOf(response)
				if v.MethodByName("RawJSON").Call(nil)[0].String() != raw {
					t.Fatal("lost raw JSON")
				}
				variant := v.MethodByName("AsAny").Call(nil)[0].Interface()
				if strings.Contains(raw, "future_policy") {
					if variant != nil {
						t.Fatal("unknown policy became a known branch")
					}
				} else {
					policy, ok := variant.(managed.ManagedAgentsAlwaysDenyPolicy)
					if !ok || policy.Type != "always_deny" || policy.JSON.ExtraFields["future"].Raw() != "7" {
						t.Fatal(variant)
					}
					parityJSON(t, policy.ToParam(), raw)
				}
			}
		})
	}
	parityJSON(t, managed.ManagedAgentsAlwaysDenyPolicyParam{}, `{"type":"always_deny"}`)
	parityJSON(t, managed.ManagedAgentsBashToolConfigParams{}, `{"name":"bash"}`)
	parityJSON(t, param.NullStruct[managed.ManagedAgentsBashToolConfigParamsPermissionPolicyUnion](), `null`)
	_, err := json.Marshal(managed.ManagedAgentsBashToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny, OfAlwaysAsk: &managed.ManagedAgentsAlwaysAskPolicyParam{Type: "always_ask"}})
	if err == nil {
		t.Fatal("multiple union branches accepted")
	}
}

func TestCapabilityParityMediaToolConfigs(t *testing.T) {
	policy := managed.ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: &managed.ManagedAgentsAlwaysDenyPolicyParam{Type: "always_deny"}}
	configs := []managed.ManagedAgentsAgentToolConfigParamsUnion{
		{OfImageGen: &managed.ManagedAgentsImageGenToolConfigParams{Name: "ImageGen", Enabled: managed.Bool(false), PermissionPolicy: policy}},
		{OfImageSearch: &managed.ManagedAgentsImageSearchToolConfigParams{Name: "ImageSearch", Enabled: managed.Bool(false), PermissionPolicy: policy}},
		{OfVideoGen: &managed.ManagedAgentsVideoGenToolConfigParams{Name: "VideoGen", Enabled: managed.Bool(false), PermissionPolicy: policy}},
		{OfDeliverArtifacts: &managed.ManagedAgentsDeliverArtifactsToolConfigParams{Name: "DeliverArtifacts", Enabled: managed.Bool(false), PermissionPolicy: policy}},
	}
	for i, name := range []string{"ImageGen", "ImageSearch", "VideoGen", "DeliverArtifacts"} {
		t.Run(name, func(t *testing.T) {
			want := `{"name":"` + name + `","enabled":false,"permission_policy":{"type":"always_deny"}}`
			parityJSON(t, configs[i], want)
			var decoded managed.ManagedAgentsAgentToolConfigParamsUnion
			if err := json.Unmarshal([]byte(want), &decoded); err != nil {
				t.Fatal(err)
			}
			if decoded.GetName() == nil || *decoded.GetName() != name || decoded.GetEnabled() == nil || *decoded.GetEnabled() || decoded.GetType() != nil {
				t.Fatal(decoded)
			}
			if _, ok := decoded.GetPermissionPolicy().AsAny().(*managed.ManagedAgentsAlwaysDenyPolicyParam); !ok {
				t.Fatal("missing request deny helper")
			}
			if typ := decoded.GetPermissionPolicy().GetType(); typ == nil || *typ != "always_deny" {
				t.Fatal(typ)
			}
			parityJSON(t, decoded, want)
			var nested managed.ManagedAgentsAgentToolset20260401Params
			if err := apijson.UnmarshalRoot([]byte(`{"type":"agent_toolset_20260401","configs":[`+want+`]}`), &nested); err != nil {
				t.Fatal(err)
			}
			if nested.Configs[0].GetName() == nil || *nested.Configs[0].GetName() != name {
				t.Fatal("nested decoder lost canonical name")
			}
			var response managed.ManagedAgentsAgentToolConfigUnion
			if err := json.Unmarshal([]byte(want), &response); err != nil {
				t.Fatal(err)
			}
			if response.JSON.Enabled.Valid() == false || response.AsAny() == nil {
				t.Fatal(response)
			}
			v := reflect.ValueOf(response.AsAny())
			if v.FieldByName("Name").String() != name {
				t.Fatal(v)
			}
			resolved := v.FieldByName("PermissionPolicy").Interface().(managed.ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion)
			if _, ok := resolved.AsAny().(managed.ManagedAgentsAlwaysDenyPolicy); !ok {
				t.Fatal(resolved)
			}
			parityJSON(t, v.MethodByName("ToParam").Call(nil)[0].Interface(), want)
		})
	}
	for i, request := range []any{managed.ManagedAgentsImageGenToolConfigParams{}, managed.ManagedAgentsImageSearchToolConfigParams{}, managed.ManagedAgentsVideoGenToolConfigParams{}, managed.ManagedAgentsDeliverArtifactsToolConfigParams{}} {
		name := []string{"ImageGen", "ImageSearch", "VideoGen", "DeliverArtifacts"}[i]
		parityJSON(t, request, `{"name":"`+name+`"}`)
	}
	var unknown managed.ManagedAgentsAgentToolConfigUnion
	raw := `{"name":"FutureTool","future":7}`
	if err := json.Unmarshal([]byte(raw), &unknown); err != nil {
		t.Fatal(err)
	}
	if unknown.AsAny() != nil || unknown.RawJSON() != raw {
		t.Fatal("lost unknown tool")
	}
}

func TestCapabilityParityMediaAliases(t *testing.T) {
	for canonical, aliases := range map[string][]string{
		"ImageGen":         {"image_gen", "imagegen", "IMAGE_GEN"},
		"ImageSearch":      {"image_search", "imagesearch", "IMAGE_SEARCH"},
		"VideoGen":         {"video_gen", "videogen", "VIDEO_GEN"},
		"DeliverArtifacts": {"deliver_artifacts", "deliverartifacts", "DELIVER_ARTIFACTS"},
	} {
		for _, name := range aliases {
			t.Run(name, func(t *testing.T) {
				raw := `{"name":"` + name + `","enabled":false,"permission_policy":{"type":"always_deny"}}`
				var request managed.ManagedAgentsAgentToolConfigParamsUnion
				if err := json.Unmarshal([]byte(raw), &request); err != nil {
					t.Fatal(err)
				}
				if reflect.ValueOf(request).FieldByName("Of" + canonical).IsNil() {
					t.Fatal("alias selected the wrong tool")
				}
				parityJSON(t, request, raw)
				var toolset managed.ManagedAgentsAgentToolset20260401Params
				if err := apijson.UnmarshalRoot([]byte(`{"type":"agent_toolset_20260401","configs":[`+raw+`]}`), &toolset); err != nil {
					t.Fatal(err)
				}
				parityJSON(t, toolset.Configs[0], raw)
				var response managed.ManagedAgentsAgentToolConfigUnion
				if err := json.Unmarshal([]byte(raw), &response); err != nil {
					t.Fatal(err)
				}
				if response.AsAny() == nil || reflect.TypeOf(response.AsAny()).Name() != "ManagedAgents"+canonical+"ToolConfig" {
					t.Fatal("response alias selected the wrong tool")
				}
				value := reflect.ValueOf(response.AsAny())
				parityJSON(t, value.MethodByName("ToParam").Call(nil)[0].Interface(), raw)
			})
		}
	}
}

func TestCapabilityParityBrowserUnions(t *testing.T) {
	browser := &managed.ManagedBrowserToolset20260714Param{}
	want := `{"type":"browser_toolset_20260714"}`
	parityJSON(t, browser, want)
	for _, request := range []any{
		managed.AgentNewParamsToolUnion{OfBrowserToolset20260714: browser},
		managed.AgentUpdateParamsToolUnion{OfBrowserToolset20260714: browser},
		managed.ManagedAgentsAgentWithOverridesParamsToolUnion{OfBrowserToolset20260714: browser},
		managed.ManagedAgentsSessionAgentUpdateToolUnionParam{OfBrowserToolset20260714: browser},
	} {
		parityJSON(t, request, want)
		decoded := reflect.New(reflect.TypeOf(request))
		if err := json.Unmarshal([]byte(want), decoded.Interface()); err != nil {
			t.Fatal(err)
		}
		if decoded.Elem().FieldByName("OfBrowserToolset20260714").IsNil() {
			t.Fatal("missing browser discriminator")
		}
		parityJSON(t, decoded.Interface(), want)
		getType := decoded.MethodByName("GetType").Call(nil)[0].Interface().(*string)
		if getType == nil || *getType != "browser_toolset_20260714" {
			t.Fatal(getType)
		}
		if decoded.MethodByName("GetConfigs").Call(nil)[0].MethodByName("AsAny").Call(nil)[0].Interface() != nil {
			t.Fatal("Browser must not expose configs")
		}
	}
	for _, response := range []any{&managed.ManagedAgentsAgentToolUnion{}, &managed.ManagedAgentsSessionAgentToolUnion{}, &managed.ManagedAgentsSessionThreadAgentToolUnion{}} {
		if err := json.Unmarshal([]byte(want), response); err != nil {
			t.Fatal(err)
		}
		variant := reflect.ValueOf(response).MethodByName("AsAny").Call(nil)[0].Interface()
		b, ok := variant.(managed.ManagedBrowserToolset20260714)
		if !ok || string(b.Type) != "browser_toolset_20260714" || !b.JSON.Type.Valid() {
			t.Fatal(variant)
		}
		parityJSON(t, b.ToParam(), want)
	}
	if reflect.TypeOf(managed.ManagedBrowserToolset20260714{}).NumField() != 2 || reflect.TypeOf(*browser).NumField() != 2 {
		t.Fatal("Browser has extra configuration fields")
	}
}

func TestCapabilityParityManagedHTTPRoundtrip(t *testing.T) {
	deny := &managed.ManagedAgentsAlwaysDenyPolicyParam{Type: "always_deny"}
	toolset := &managed.ManagedAgentsAgentToolset20260401Params{Type: "agent_toolset_20260401", Configs: []managed.ManagedAgentsAgentToolConfigParamsUnion{
		{OfBash: &managed.ManagedAgentsBashToolConfigParams{Name: constant.Bash("Bash"), PermissionPolicy: managed.ManagedAgentsBashToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfEdit: &managed.ManagedAgentsEditToolConfigParams{Name: constant.Edit("Edit"), PermissionPolicy: managed.ManagedAgentsEditToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfRead: &managed.ManagedAgentsReadToolConfigParams{Name: constant.Read("Read"), PermissionPolicy: managed.ManagedAgentsReadToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfWrite: &managed.ManagedAgentsWriteToolConfigParams{Name: constant.Write("Write"), PermissionPolicy: managed.ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfGlob: &managed.ManagedAgentsGlobToolConfigParams{Name: constant.Glob("Glob"), PermissionPolicy: managed.ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfGrep: &managed.ManagedAgentsGrepToolConfigParams{Name: constant.Grep("Grep"), PermissionPolicy: managed.ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfWebFetch: &managed.ManagedAgentsWebFetchToolConfigParams{Name: constant.WebFetch("WebFetch"), PermissionPolicy: managed.ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfWebSearch: &managed.ManagedAgentsWebSearchToolConfigParams{Name: constant.WebSearch("WebSearch"), PermissionPolicy: managed.ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfImageGen: &managed.ManagedAgentsImageGenToolConfigParams{Name: "ImageGen", PermissionPolicy: managed.ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfImageSearch: &managed.ManagedAgentsImageSearchToolConfigParams{Name: "ImageSearch", PermissionPolicy: managed.ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfVideoGen: &managed.ManagedAgentsVideoGenToolConfigParams{Name: "VideoGen", PermissionPolicy: managed.ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
		{OfDeliverArtifacts: &managed.ManagedAgentsDeliverArtifactsToolConfigParams{Name: "DeliverArtifacts", PermissionPolicy: managed.ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}},
	}, DefaultConfig: managed.ManagedAgentsAgentToolsetDefaultConfigParams{PermissionPolicy: managed.ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}}
	mcp := &managed.ManagedAgentsMCPToolsetParams{Type: "mcp_toolset", MCPServerName: "example", Configs: []managed.ManagedAgentsMCPToolConfigParams{{Name: "forecast", PermissionPolicy: managed.ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}}, DefaultConfig: managed.ManagedAgentsMCPToolsetDefaultConfigParams{PermissionPolicy: managed.ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: deny}}}
	browser := &managed.ManagedBrowserToolset20260714Param{}
	calls := 0
	client := managed.NewClient(option.WithPAT("offline-token"), option.WithBaseURL("https://managed.test/api/v1/cloud"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != "POST" || r.Header.Get("Authorization") != "Bearer offline-token" || r.Header.Get("X-Qoder-Beta") != "browser-use-2026-07-14" {
			t.Fatal(r.Method, r.Header)
		}
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		tools := body["tools"]
		if strings.Contains(r.URL.Path, "/sessions") {
			var agent map[string]json.RawMessage
			if err := json.Unmarshal(body["agent"], &agent); err != nil {
				t.Fatal(err)
			}
			tools = agent["tools"]
			if calls == 3 && !strings.Contains(string(body["metadata"]), "__qca_sdk_created") {
				t.Fatal("missing SDK metadata")
			}
		}
		var entries []map[string]json.RawMessage
		if err := json.Unmarshal(tools, &entries); err != nil {
			t.Fatal(err)
		}
		if len(entries) != 3 {
			t.Fatal(len(entries))
		}
		var decodedTools []managed.AgentNewParamsToolUnion
		if err := json.Unmarshal(tools, &decodedTools); err != nil {
			t.Fatal(err)
		}
		if decodedTools[2].OfBrowserToolset20260714 == nil {
			t.Fatal("nested decoder lost browser")
		}
		if len(decodedTools[0].OfAgentToolset20260401.Configs) != 12 {
			t.Fatal("nested decoder lost configs")
		}
		for _, config := range decodedTools[0].OfAgentToolset20260401.Configs {
			if config.GetName() == nil || config.GetPermissionPolicy().GetType() == nil || *config.GetPermissionPolicy().GetType() != "always_deny" {
				t.Fatal("nested named config decoder lost deny", config)
			}
		}
		parityJSON(t, decodedTools, string(tools))
		parityJSON(t, entries[2], `{"type":"browser_toolset_20260714"}`)
		var configs []map[string]json.RawMessage
		if err := json.Unmarshal(entries[0]["configs"], &configs); err != nil {
			t.Fatal(err)
		}
		for i, name := range []string{"Bash", "Edit", "Read", "Write", "Glob", "Grep", "WebFetch", "WebSearch", "ImageGen", "ImageSearch", "VideoGen", "DeliverArtifacts"} {
			parityJSON(t, configs[i]["name"], `"`+name+`"`)
			parityJSON(t, configs[i]["permission_policy"], `{"type":"always_deny"}`)
			if _, present := configs[i]["type"]; present {
				t.Fatal("tool config must work by canonical name alone")
			}
		}
		parityJSON(t, entries[0]["default_config"], `{"permission_policy":{"type":"always_deny"}}`)
		parityJSON(t, entries[1]["default_config"], `{"permission_policy":{"type":"always_deny"}}`)
		parityJSON(t, entries[1]["configs"], `[{"name":"forecast","permission_policy":{"type":"always_deny"}}]`)
		if calls == 2 && r.URL.EscapedPath() != "/api/v1/cloud/agents/"+url.PathEscape("agent /?#%") {
			t.Fatal(r.URL)
		}
		if calls <= 2 {
			return reply(r, 200, `{"id":"agent","type":"agent","version":3,"tools":`+string(tools)+`}`), nil
		}
		return reply(r, 200, `{"id":"session","agent":{"id":"agent","type":"agent","version":3,"tools":`+string(tools)+`}}`), nil
	})}))
	newTools := []managed.AgentNewParamsToolUnion{{OfAgentToolset20260401: toolset}, {OfMCPToolset: mcp}, {OfBrowserToolset20260714: browser}}
	if typ := newTools[0].GetDefaultConfig().GetPermissionPolicy().GetType(); typ == nil || *typ != "always_deny" {
		t.Fatal("default policy helper lost deny")
	}
	res, err := client.Agents.New(context.Background(), managed.AgentNewParams{Name: "parity", Model: managed.ManagedAgentsModelConfigParams{ID: "ultimate"}, Tools: newTools, Betas: []managed.QoderBeta{"browser-use-2026-07-14"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.Version != 3 || !res.JSON.Version.Valid() {
		t.Fatal("numeric main version changed")
	}
	resolved := res.Tools[0].AsAgentToolset20260401()
	if _, ok := resolved.DefaultConfig.PermissionPolicy.AsAny().(managed.ManagedAgentsAlwaysDenyPolicy); !ok {
		t.Fatal("lost default deny")
	}
	for _, config := range resolved.Configs {
		if config.AsAny() == nil {
			t.Fatal("lost named tool branch", config.Name)
		}
		v := reflect.ValueOf(config.AsAny()).FieldByName("PermissionPolicy")
		variant := v.MethodByName("AsAny").Call(nil)[0].Interface()
		if _, ok := variant.(managed.ManagedAgentsAlwaysDenyPolicy); !ok {
			t.Fatal("lost tool deny", config.Name)
		}
	}
	if _, ok := res.Tools[1].AsMCPToolset().Configs[0].PermissionPolicy.AsAny().(managed.ManagedAgentsAlwaysDenyPolicy); !ok {
		t.Fatal("lost MCP deny")
	}
	if _, ok := res.Tools[2].AsAny().(managed.ManagedBrowserToolset20260714); !ok {
		t.Fatal("lost browser")
	}
	_, err = client.Agents.Update(context.Background(), "agent /?#%", managed.AgentUpdateParams{Tools: []managed.AgentUpdateParamsToolUnion{{OfAgentToolset20260401: toolset}, {OfMCPToolset: mcp}, {OfBrowserToolset20260714: browser}}}, option.WithHeader("X-Qoder-Beta", "browser-use-2026-07-14"))
	if err != nil {
		t.Fatal(err)
	}
	session, err := client.Sessions.New(context.Background(), managed.SessionNewParams{EnvironmentID: "env", Agent: managed.SessionNewParamsAgentUnion{OfManagedAgentsAgentWithOverridess: &managed.ManagedAgentsAgentWithOverridesParams{ID: "agent", Type: "agent_with_overrides", Version: managed.Int(3), Tools: []managed.ManagedAgentsAgentWithOverridesParamsToolUnion{{OfAgentToolset20260401: toolset}, {OfMCPToolset: mcp}, {OfBrowserToolset20260714: browser}}}}}, option.WithHeader("X-Qoder-Beta", "browser-use-2026-07-14"))
	if err != nil {
		t.Fatal(err)
	}
	if session.Agent.Version != 3 || !session.Agent.JSON.Version.Valid() || len(session.Agent.Tools[0].AsAgentToolset20260401().Configs) != 12 {
		t.Fatal(session.Agent)
	}
	if _, ok := session.Agent.Tools[2].AsAny().(managed.ManagedBrowserToolset20260714); !ok {
		t.Fatal("lost session browser")
	}
	_, err = client.Sessions.Update(context.Background(), "session", managed.SessionUpdateParams{Agent: managed.ManagedAgentsSessionAgentUpdateParam{Tools: []managed.ManagedAgentsSessionAgentUpdateToolUnionParam{{OfAgentToolset20260401: toolset}, {OfMCPToolset: mcp}, {OfBrowserToolset20260714: browser}}}}, option.WithHeader("X-Qoder-Beta", "browser-use-2026-07-14"))
	if err != nil || calls != 4 {
		t.Fatal(err, calls)
	}
}

func TestCapabilityParityTopMCPPolicies(t *testing.T) {
	for _, policyType := range []string{"always_allow", "always_ask", "always_deny"} {
		t.Run(policyType, func(t *testing.T) {
			var policy managed.ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion
			if err := json.Unmarshal([]byte(`{"type":"`+policyType+`"}`), &policy); err != nil {
				t.Fatal(err)
			}
			mcp := &managed.ManagedAgentsMCPToolsetParams{Type: "mcp_toolset", MCPServerName: "example", PermissionPolicy: policy}
			raw := `{"type":"mcp_toolset","mcp_server_name":"example","permission_policy":{"type":"` + policyType + `"}}`
			parityJSON(t, mcp, raw)
			for _, request := range []any{
				managed.AgentNewParamsToolUnion{OfMCPToolset: mcp},
				managed.AgentUpdateParamsToolUnion{OfMCPToolset: mcp},
				managed.ManagedAgentsAgentWithOverridesParamsToolUnion{OfMCPToolset: mcp},
				managed.ManagedAgentsSessionAgentUpdateToolUnionParam{OfMCPToolset: mcp},
			} {
				t.Run(reflect.TypeOf(request).Name(), func(t *testing.T) {
					parityJSON(t, request, raw)
					decoded := reflect.New(reflect.TypeOf(request))
					if err := json.Unmarshal([]byte(raw), decoded.Interface()); err != nil {
						t.Fatal(err)
					}
					parityJSON(t, decoded.Interface(), raw)
					top := decoded.MethodByName("GetPermissionPolicy").Call(nil)[0].Interface().(managed.ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion)
					if typ := top.GetType(); typ == nil || *typ != policyType || top.AsAny() == nil {
						t.Fatal("top-level policy helper lost its branch", top)
					}
					parityJSON(t, top.AsAny(), `{"type":"`+policyType+`"}`)
				})
			}
			for _, response := range []any{managed.ManagedAgentsAgentToolUnion{}, managed.ManagedAgentsSessionAgentToolUnion{}, managed.ManagedAgentsSessionThreadAgentToolUnion{}} {
				t.Run(reflect.TypeOf(response).Name(), func(t *testing.T) {
					decoded := reflect.New(reflect.TypeOf(response))
					if err := json.Unmarshal([]byte(raw), decoded.Interface()); err != nil {
						t.Fatal(err)
					}
					top := decoded.Elem().FieldByName("PermissionPolicy").Interface().(managed.ManagedAgentsMCPToolConfigPermissionPolicyUnion)
					if top.Type != policyType || top.AsAny() == nil {
						t.Fatal("implicit union lost top-level policy", top)
					}
					resolved := decoded.MethodByName("AsMCPToolset").Call(nil)[0].Interface().(managed.ManagedAgentsMCPToolset)
					if resolved.PermissionPolicy.Type != policyType || !resolved.JSON.PermissionPolicy.Valid() || resolved.RawJSON() != raw {
						t.Fatal("MCP response lost top-level policy", resolved)
					}
				})
			}
		})
	}
	legacy := &managed.ManagedAgentsMCPToolsetParams{Type: "mcp_toolset", MCPServerName: "example", DefaultConfig: managed.ManagedAgentsMCPToolsetDefaultConfigParams{PermissionPolicy: managed.ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion{OfAlwaysDeny: &managed.ManagedAgentsAlwaysDenyPolicyParam{}}}}
	parityJSON(t, legacy, `{"type":"mcp_toolset","mcp_server_name":"example","default_config":{"permission_policy":{"type":"always_deny"}}}`)
	if (managed.AgentNewParamsToolUnion{OfMCPToolset: legacy}).GetPermissionPolicy().AsAny() != nil {
		t.Fatal("legacy policy was silently promoted to the top level")
	}
	legacy.PermissionPolicy = managed.ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion{OfAlwaysAsk: &managed.ManagedAgentsAlwaysAskPolicyParam{Type: "always_ask"}}
	parityJSON(t, legacy, `{"type":"mcp_toolset","mcp_server_name":"example","permission_policy":{"type":"always_ask"},"default_config":{"permission_policy":{"type":"always_deny"}}}`)
	legacy.PermissionPolicy = param.NullStruct[managed.ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion]()
	parityJSON(t, legacy, `{"type":"mcp_toolset","mcp_server_name":"example","permission_policy":null,"default_config":{"permission_policy":{"type":"always_deny"}}}`)
	for _, request := range []any{managed.AgentNewParamsToolUnion{}, managed.AgentUpdateParamsToolUnion{}, managed.ManagedAgentsAgentWithOverridesParamsToolUnion{}, managed.ManagedAgentsSessionAgentUpdateToolUnionParam{}} {
		policy := reflect.ValueOf(request).MethodByName("GetPermissionPolicy").Call(nil)[0].Interface().(managed.ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion)
		if policy.AsAny() != nil || policy.GetType() != nil || !param.IsOmitted(policy) {
			t.Fatal("omitted policy acquired a default")
		}
	}
	for _, response := range []any{managed.ManagedAgentsMCPToolset{}, managed.ManagedAgentsAgentToolUnion{}, managed.ManagedAgentsSessionAgentToolUnion{}, managed.ManagedAgentsSessionThreadAgentToolUnion{}} {
		for _, policy := range []string{"", `null`, `{"type":"future_policy","future":7}`} {
			raw := `{"type":"mcp_toolset","mcp_server_name":"example"`
			if policy != "" {
				raw += `,"permission_policy":` + policy
			}
			raw += `}`
			decoded := reflect.New(reflect.TypeOf(response))
			if err := json.Unmarshal([]byte(raw), decoded.Interface()); err != nil {
				t.Fatal(err)
			}
			top := decoded.Elem().FieldByName("PermissionPolicy").Interface().(managed.ManagedAgentsMCPToolConfigPermissionPolicyUnion)
			meta := decoded.Elem().FieldByName("JSON").FieldByName("PermissionPolicy").Interface().(interface{ Raw() string })
			if meta.Raw() != policy || top.AsAny() != nil {
				t.Fatal("omitted, null or unknown policy became a known default", top, meta.Raw())
			}
			if policy != "" && policy != "null" && top.RawJSON() != policy {
				t.Fatal("lost unknown policy JSON")
			}
		}
	}
}

// These contracts inspect HTTP bodies and responses, not real MCP tool execution.
func TestCapabilityParityTopMCPDefaultDenyAPIContracts(t *testing.T) {
	mcp := &managed.ManagedAgentsMCPToolsetParams{Type: "mcp_toolset", MCPServerName: "example", PermissionPolicy: managed.ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion{OfAlwaysDeny: &managed.ManagedAgentsAlwaysDenyPolicyParam{}}}
	servers := []managed.ManagedAgentsURLMCPServerParams{{Type: "url", Name: "example", URL: "https://mcp.example.invalid"}}
	toolJSON := `{"type":"mcp_toolset","mcp_server_name":"example","permission_policy":{"type":"always_deny"}}`
	for _, operation := range []string{"create agent", "update agent", "agent overrides", "session agent patch"} {
		t.Run(operation, func(t *testing.T) {
			path := map[string]string{"create agent": "/agents", "update agent": "/agents/agent", "agent overrides": "/sessions", "session agent patch": "/sessions/session"}[operation]
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Method != http.MethodPost || r.URL.Path != "/api/v1/cloud"+path || r.Header.Get("Authorization") != "Bearer offline-token" {
					t.Fatal(r.Method, r.URL)
				}
				var body map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				agent := body
				if strings.Contains(path, "/sessions") {
					if err := json.Unmarshal(body["agent"], &agent); err != nil {
						t.Fatal(err)
					}
					if operation == "agent overrides" && !strings.Contains(string(body["metadata"]), "__qca_sdk_created") {
						t.Fatal("missing SDK metadata")
					}
				}
				var tools []map[string]json.RawMessage
				if err := json.Unmarshal(agent["tools"], &tools); err != nil {
					t.Fatal(err)
				}
				if len(tools) != 1 {
					t.Fatal(tools)
				}
				parityJSON(t, tools[0], toolJSON)
				if _, present := tools[0]["configs"]; present {
					t.Fatal("default-deny contract must not depend on configs")
				}
				if _, present := tools[0]["default_config"]; present {
					t.Fatal("top-level deny was moved to default_config")
				}
				parityJSON(t, agent["mcp_servers"], `[{"type":"url","name":"example","url":"https://mcp.example.invalid"}]`)
				if strings.Contains(path, "/sessions") {
					return reply(r, 200, `{"id":"session","agent":{"id":"agent","type":"agent","version":3,"tools":[`+toolJSON+`]}}`), nil
				}
				return reply(r, 200, `{"id":"agent","type":"agent","version":3,"tools":[`+toolJSON+`]}`), nil
			}, option.WithPAT("offline-token"))
			var top managed.ManagedAgentsMCPToolConfigPermissionPolicyUnion
			var resolved managed.ManagedAgentsMCPToolset
			switch operation {
			case "create agent":
				res, err := client.Agents.New(context.Background(), managed.AgentNewParams{Name: "default-deny", Model: managed.ManagedAgentsModelConfigParams{ID: "ultimate"}, MCPServers: servers, Tools: []managed.AgentNewParamsToolUnion{{OfMCPToolset: mcp}}})
				if err != nil {
					t.Fatal(err)
				}
				top, resolved = res.Tools[0].PermissionPolicy, res.Tools[0].AsMCPToolset()
			case "update agent":
				res, err := client.Agents.Update(context.Background(), "agent", managed.AgentUpdateParams{MCPServers: servers, Tools: []managed.AgentUpdateParamsToolUnion{{OfMCPToolset: mcp}}})
				if err != nil {
					t.Fatal(err)
				}
				top, resolved = res.Tools[0].PermissionPolicy, res.Tools[0].AsMCPToolset()
			case "agent overrides":
				res, err := client.Sessions.New(context.Background(), managed.SessionNewParams{EnvironmentID: "env", Agent: managed.SessionNewParamsAgentUnion{OfManagedAgentsAgentWithOverridess: &managed.ManagedAgentsAgentWithOverridesParams{ID: "agent", Type: "agent_with_overrides", Version: managed.Int(3), MCPServers: servers, Tools: []managed.ManagedAgentsAgentWithOverridesParamsToolUnion{{OfMCPToolset: mcp}}}}})
				if err != nil {
					t.Fatal(err)
				}
				top, resolved = res.Agent.Tools[0].PermissionPolicy, res.Agent.Tools[0].AsMCPToolset()
			case "session agent patch":
				res, err := client.Sessions.Update(context.Background(), "session", managed.SessionUpdateParams{Agent: managed.ManagedAgentsSessionAgentUpdateParam{MCPServers: servers, Tools: []managed.ManagedAgentsSessionAgentUpdateToolUnionParam{{OfMCPToolset: mcp}}}})
				if err != nil {
					t.Fatal(err)
				}
				top, resolved = res.Agent.Tools[0].PermissionPolicy, res.Agent.Tools[0].AsMCPToolset()
			}
			if calls != 1 || top.Type != "always_deny" || !resolved.JSON.PermissionPolicy.Valid() {
				t.Fatal("lost top-level default deny", calls, top, resolved)
			}
			if _, ok := resolved.PermissionPolicy.AsAny().(managed.ManagedAgentsAlwaysDenyPolicy); !ok {
				t.Fatal("lost deny response branch")
			}
		})
	}
}
