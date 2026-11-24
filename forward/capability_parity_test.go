package forward_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
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

func TestCapabilityParitySessionMultiagent(t *testing.T) {
	roster := forward.MultiagentConfigParam{Type: "coordinator", Agents: []forward.MultiagentEntryParam{
		{OfString: forward.String("tmpl_child")}, {Type: "agent", TemplateID: forward.String("tmpl_second")}, {Type: "agent", ID: forward.String("agent_child")}, {Type: "self"}, {Type: "advisor", Model: forward.String("ultimate")},
	}}
	for _, tc := range []struct {
		name   string
		config forward.SessionUpdateParamsConfigParam
		want   string
	}{
		{"omitted", forward.SessionUpdateParamsConfigParam{}, `{}`},
		{"null", forward.SessionUpdateParamsConfigParam{Multiagent: param.NullStruct[forward.MultiagentConfigParam]()}, `{"config":{"multiagent":null}}`},
		{"roster", forward.SessionUpdateParamsConfigParam{Multiagent: roster}, `{"config":{"multiagent":{"type":"coordinator","agents":["tmpl_child",{"type":"agent","template_id":"tmpl_second"},{"type":"agent","id":"agent_child"},{"type":"self"},{"type":"advisor","model":"ultimate"}]}}}`},
		{"empty", forward.SessionUpdateParamsConfigParam{Multiagent: forward.MultiagentConfigParam{Type: "coordinator", Agents: []forward.MultiagentEntryParam{}}}, `{"config":{"multiagent":{"type":"coordinator","agents":[]}}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := forward.SessionUpdateParams{Config: tc.config}
			parityJSON(t, p, tc.want)
			var decoded forward.SessionUpdateParams
			if err := json.Unmarshal([]byte(tc.want), &decoded); err != nil {
				t.Fatal(err)
			}
			parityJSON(t, decoded, tc.want)
			client := testClient(func(r *http.Request) (*http.Response, error) {
				if r.Method != "POST" || r.URL.EscapedPath() != "/api/v1/forward/sessions/"+url.PathEscape(pathSegment) {
					t.Fatal(r.URL)
				}
				var body any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				parityJSON(t, body, tc.want)
				if tc.name == "empty" {
					return reply(r, 400, `{"error":{"type":"invalid_request","message":"multiagent.agents must be a non-empty array"}}`), nil
				}
				return reply(r, 200, tc.want), nil
			})
			res, err := client.Sessions.Update(context.Background(), pathSegment, p)
			if tc.name == "empty" {
				if err == nil {
					t.Fatal("expected the service to reject an empty roster")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tc.name == "null" && res.Config.JSON.Multiagent.Raw() != "null" {
				t.Fatal("lost explicit null")
			}
			if tc.name == "roster" && (res.Config.Multiagent.Agents[0].OfString != "tmpl_child" || res.Config.Multiagent.Agents[1].TemplateID != "tmpl_second" || res.Config.Multiagent.Agents[2].ID != "agent_child") {
				t.Fatal(res.Config.Multiagent)
			}
		})
	}
	if _, ok := reflect.TypeOf(forward.SessionNewParamsConfigParam{}).FieldByName("Multiagent"); ok {
		t.Fatal("create config must not offer multiagent")
	}
	var snapshot forward.Session
	if err := json.Unmarshal([]byte(`{"template":{"id":"tmpl_main","version":3},"config":{"multiagent":null}}`), &snapshot); err != nil {
		t.Fatal(err)
	}
	if snapshot.Template.Version != 3 || !snapshot.Template.JSON.Version.Valid() {
		t.Fatal("main agent version changed")
	}
}

func TestCapabilityParityAutoTitle(t *testing.T) {
	for _, tc := range []struct {
		name    string
		auto    param.Opt[bool]
		title   param.Opt[string]
		present bool
		bad     bool
	}{
		{"omitted", param.Opt[bool]{}, param.Opt[string]{}, false, false},
		{"false", forward.Bool(false), param.Opt[string]{}, true, false},
		{"true", forward.Bool(true), param.Opt[string]{}, true, false},
		{"title", forward.Bool(true), forward.String("title"), true, true},
		{"empty title", forward.Bool(true), forward.String(""), true, true},
		{"null title", forward.Bool(true), param.Null[string](), true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				var body map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				_, present := body["auto_title"]
				if present != tc.present {
					t.Fatal(body)
				}
				if present {
					parityJSON(t, body["auto_title"], map[bool]string{true: "true", false: "false"}[tc.auto.Value])
				}
				if _, present := body["title"]; present != tc.bad {
					t.Fatal(body)
				}
				if !strings.Contains(string(body["metadata"]), "__qca_sdk_created") {
					t.Fatal("missing metadata")
				}
				if tc.bad {
					return reply(r, 400, `{"error":{"type":"invalid_request_error","message":"title must not be supplied"}}`), nil
				}
				return reply(r, 200, `{"id":"sess_auto","title":"generated"}`), nil
			})
			_, err := client.Sessions.New(context.Background(), forward.SessionNewParams{IdentityID: "idn", TemplateID: "tmpl", AutoTitle: tc.auto, Title: tc.title})
			if calls != 1 || (err != nil) != tc.bad {
				t.Fatalf("calls=%d err=%v", calls, err)
			}
		})
	}
	for _, typ := range []reflect.Type{reflect.TypeOf(forward.SessionUpdateParams{}), reflect.TypeOf(forward.Session{})} {
		if _, ok := typ.FieldByName("AutoTitle"); ok {
			t.Fatal("auto_title is create-only")
		}
	}
}

func TestCapabilityParityTemplateA2A(t *testing.T) {
	for _, tc := range []struct {
		name string
		a2a  forward.TemplateA2AConfigParam
		want string
	}{
		{"omitted", forward.TemplateA2AConfigParam{}, `{}`},
		{"empty", param.Override[forward.TemplateA2AConfigParam](json.RawMessage(`{}`)), `{"a2a":{}}`},
		{"clear", param.NullStruct[forward.TemplateA2AConfigParam](), `{"a2a":null}`},
		{"disable", forward.TemplateA2AConfigParam{Enabled: forward.Bool(false)}, `{"a2a":{"enabled":false}}`},
		{"clear description", forward.TemplateA2AConfigParam{Description: forward.String("")}, `{"a2a":{"description":""}}`},
		{"clear skills", forward.TemplateA2AConfigParam{Skills: []forward.TemplateA2ASkillParam{}}, `{"a2a":{"skills":[]}}`},
		{"replace skills", forward.TemplateA2AConfigParam{Skills: []forward.TemplateA2ASkillParam{{ID: "summary", Name: "Summary", Description: "Public", Tags: []string{}, Examples: []string{"Summarize"}}}}, `{"a2a":{"skills":[{"id":"summary","name":"Summary","description":"Public","tags":[],"examples":["Summarize"]}]}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := forward.TemplateUpdateParams{A2A: tc.a2a}
			parityJSON(t, p, tc.want)
			var decoded forward.TemplateUpdateParams
			if err := json.Unmarshal([]byte(tc.want), &decoded); err != nil {
				t.Fatal(err)
			}
			parityJSON(t, decoded, tc.want)
			client := testClient(func(r *http.Request) (*http.Response, error) {
				var body any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				parityJSON(t, body, tc.want)
				return reply(r, 200, tc.want), nil
			})
			res, err := client.Templates.Update(context.Background(), "tmpl", p)
			if err != nil {
				t.Fatal(err)
			}
			if tc.name == "clear" && res.JSON.A2A.Raw() != "null" {
				t.Fatal("lost null")
			}
			if tc.name == "replace skills" && (res.A2A.Skills[0].ID != "summary" || res.A2A.Skills[0].Examples[0] != "Summarize") {
				t.Fatal(res.A2A)
			}
		})
	}
	client := testClient(func(r *http.Request) (*http.Response, error) {
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		parityJSON(t, body["a2a"], `{"enabled":true,"description":"Public","skills":[]}`)
		return reply(r, 200, `{"a2a":{"enabled":true,"description":"Public","skills":[]}}`), nil
	})
	res, err := client.Templates.New(context.Background(), forward.TemplateNewParams{Name: "a2a", Model: forward.ModelConfigUnionParam{OfString: forward.String("ultimate")}, EnvironmentID: "env", A2A: forward.TemplateA2AConfigParam{Enabled: forward.Bool(true), Description: forward.String("Public"), Skills: []forward.TemplateA2ASkillParam{}}})
	if err != nil || !res.A2A.Enabled {
		t.Fatal(res, err)
	}
	for _, patch := range []string{`{"enabled":null}`, `{"description":null}`, `{"skills":null}`, `{"skills":[{"id":null,"name":"Summary","description":"Public"}]}`} {
		params := forward.TemplateUpdateParams{A2A: param.Override[forward.TemplateA2AConfigParam](json.RawMessage(patch))}
		calls := 0
		client := testClient(func(r *http.Request) (*http.Response, error) {
			calls++
			var body any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			parityJSON(t, body, `{"a2a":`+patch+`}`)
			return reply(r, 400, `{"error":{"type":"invalid_request_error","message":"A2A members cannot be null"}}`), nil
		})
		_, err := client.Templates.Update(context.Background(), "tmpl", params)
		if err == nil || calls != 1 {
			t.Fatal("member-null validation must remain server-owned", err, calls)
		}
	}
}

func TestCapabilityParityAgentCardAPIContracts(t *testing.T) {
	for _, status := range []int{200, 400} {
		client := contractClient(t, "getTemplateAgentCard", status)
		res, err := client.Templates.GetAgentCard(context.Background(), pathSegment, contractParams[forward.TemplateGetAgentCardParams](t, "getTemplateAgentCard"))
		checkError(t, status, err)
		if status == 200 {
			checkDecoded(t, res)
		}
	}
	for _, mode := range []string{"pat", "sa"} {
		for _, version := range []string{"", "0.3", "1.0"} {
			t.Run(mode+version, func(t *testing.T) {
				params := forward.TemplateGetAgentCardParams{}
				if mode == "pat" {
					params.IdentityID = forward.String(pathSegment)
				}
				if version != "" {
					params.A2AVersion = forward.String(version)
				}
				calls := 0
				client := testClient(func(r *http.Request) (*http.Response, error) {
					calls++
					if r.Method != "GET" || r.URL.EscapedPath() != "/api/v1/forward/templates/"+url.PathEscape(pathSegment)+"/a2a/agent-card.json" {
						t.Fatal(r.URL)
					}
					if r.Header.Get("Authorization") != "Bearer card-token" || r.Header.Get("A2A-Version") != version || r.Header.Get("X-Test") != "request" {
						t.Fatal(r.Header)
					}
					q := url.Values{}
					if mode == "pat" {
						q.Set("identity_id", pathSegment)
					}
					if r.URL.RawQuery != q.Encode() {
						t.Fatal(r.URL.RawQuery)
					}
					body := `{"name":"Public","description":"Summary","version":"1.0.0","skills":[{"id":"summary","name":"Summary","description":"Public","tags":["text"],"examples":["Summarize"]}],"defaultInputModes":["text/plain"],"defaultOutputModes":["text/plain"],"capabilities":{"streaming":true},`
					if version == "1.0" {
						body += `"supportedInterfaces":[{"url":"https://never-follow.test/a2a","protocolBinding":"JSONRPC","protocolVersion":"1.0","future":7}]}`
					} else {
						body += `"url":"https://never-follow.test/a2a"}`
					}
					return reply(r, 200, body), nil
				}, option.WithAuthMode(mode), option.WithAuthToken("card-token"), option.WithHeader("X-Test", "client"))
				client.Templates.Options = append(client.Templates.Options, option.WithHeader("X-Test", "service"))
				var raw *http.Response
				res, err := client.Templates.GetAgentCard(context.Background(), pathSegment, params, option.WithHeader("X-Test", "request"), option.WithResponseInto(&raw))
				if err != nil {
					t.Fatal(err)
				}
				if calls != 1 || raw.Header.Get("X-Request-Id") != "req_contract" {
					t.Fatal("unexpected follow-up request or missing metadata")
				}
				if res.Name != "Public" || res.Description != "Summary" || res.Version != "1.0.0" || res.Skills[0].ID != "summary" || res.DefaultInputModes[0] != "text/plain" || res.DefaultOutputModes[0] != "text/plain" || res.JSON.ExtraFields["capabilities"].Raw() != `{"streaming":true}` {
					t.Fatal(res)
				}
				if version == "1.0" && (res.SupportedInterfaces[0].ProtocolVersion != "1.0" || res.SupportedInterfaces[0].JSON.ExtraFields["future"].Raw() != "7") {
					t.Fatal(res.SupportedInterfaces)
				}
				if version != "1.0" && res.URL != "https://never-follow.test/a2a" {
					t.Fatal(res.URL)
				}
			})
		}
	}
}

func TestCapabilityParityEnvironmentStorage(t *testing.T) {
	buckets := make([]forward.EnvironmentStorageBucketParam, 7)
	for i := range buckets {
		buckets[i] = forward.EnvironmentStorageBucketParam{Name: string(rune('a' + i)), Provider: "aliyun-oss", Region: "cn-hangzhou", Bucket: "example-bucket", Access: forward.String("read_only"), Credential: forward.EnvironmentStorageCredentialParam{Type: "ram_role", RoleARN: forward.String("acs:ram::1234567890123456:role/example")}}
	}
	buckets[0].Access = forward.String("read_write")
	buckets[0].Prefix = forward.String("team/data/")
	buckets[0].AcceleratorEndpoint = forward.String("example.oss-data-acc.aliyuncs.com")
	buckets[0].Credential = forward.EnvironmentStorageCredentialParam{Type: "access_key", AccessKeyID: forward.String("test-key-id"), AccessKeySecret: forward.String("test-secret")}
	config := forward.EnvironmentConfigParams{Type: forward.String("cloud"), Packages: forward.EnvironmentConfigPackagesParam{Apt: []string{"git"}, Npm: []string{}, Pip: []string{"PyYAML"}}, SetupScript: forward.String("setup"), StartupScript: forward.String("startup"), Storage: forward.EnvironmentStorageParam{Buckets: buckets}}
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != "POST" || r.Header.Get("X-Qoder-Beta") != "oss-bucket-2026-09-23" {
			t.Fatal(r.Method, r.Header)
		}
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		var cfg map[string]json.RawMessage
		if err := json.Unmarshal(body["config"], &cfg); err != nil {
			t.Fatal(err)
		}
		var storage struct {
			Buckets []map[string]json.RawMessage `json:"buckets"`
		}
		if err := json.Unmarshal(cfg["storage"], &storage); err != nil {
			t.Fatal(err)
		}
		if len(storage.Buckets) != 7 {
			t.Fatal(len(storage.Buckets))
		}
		if calls == 1 {
			parityJSON(t, storage.Buckets[0]["credential"], `{"type":"access_key","access_key_id":"test-key-id","access_key_secret":"test-secret"}`)
			parityJSON(t, storage.Buckets[1]["credential"], `{"type":"ram_role","role_arn":"acs:ram::1234567890123456:role/example"}`)
			parityJSON(t, cfg["packages"], `{"apt":["git"],"npm":[],"pip":["PyYAML"]}`)
			parityJSON(t, cfg["setup_script"], `"setup"`)
			parityJSON(t, cfg["startup_script"], `"startup"`)
		} else {
			for _, bucket := range storage.Buckets {
				if _, present := bucket["credential"]; present {
					t.Fatal("credential must be omitted on update")
				}
			}
		}
		return reply(r, 200, `{"id":"env","config":{"type":"cloud","storage":{"buckets":[{"name":"a","provider":"aliyun-oss","region":"cn-hangzhou","bucket":"example-bucket","access":"read_write","prefix":"team/data/","accelerator_endpoint":"example.oss-data-acc.aliyuncs.com","credential":{"type":"access_key","access_key_id":"test****id"}},{"name":"b","credential":{"type":"ram_role"}}]}}}`), nil
	})
	res, err := client.Environments.New(context.Background(), forward.EnvironmentNewParams{Name: "storage", Config: config}, option.WithHeader("X-Qoder-Beta", "oss-bucket-2026-09-23"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Config.Storage.Buckets[0].Credential.AccessKeyID != "test****id" || !res.Config.Storage.Buckets[0].Credential.JSON.AccessKeyID.Valid() || res.Config.Storage.Buckets[1].Credential.Type != "ram_role" || res.Config.Storage.Buckets[1].Credential.JSON.AccessKeyID.Valid() || strings.Contains(res.RawJSON(), "test-secret") || strings.Contains(res.RawJSON(), "role_arn") {
		t.Fatal("invalid credential read view")
	}
	for i := range config.Storage.Buckets {
		config.Storage.Buckets[i].Credential = forward.EnvironmentStorageCredentialParam{}
	}
	_, err = client.Environments.Update(context.Background(), "env", forward.EnvironmentUpdateParams{Config: config}, option.WithHeader("X-Qoder-Beta", "oss-bucket-2026-09-23"))
	if err != nil || calls != 2 {
		t.Fatal(err, calls)
	}
	parityJSON(t, forward.EnvironmentUpdateParams{}, `{}`)
	parityJSON(t, forward.EnvironmentUpdateParams{Config: forward.EnvironmentConfigParams{Storage: forward.EnvironmentStorageParam{Buckets: []forward.EnvironmentStorageBucketParam{}}}}, `{"config":{"storage":{"buckets":[]}}}`)
	for _, name := range []string{"RoleARN", "RoleArn", "AccessKeySecret"} {
		if _, ok := reflect.TypeOf(forward.EnvironmentStorageCredential{}).FieldByName(name); ok {
			t.Fatal("response exposes write credential", name)
		}
	}
	var p forward.EnvironmentNewParams
	raw, err := json.Marshal(forward.EnvironmentNewParams{Name: "storage", Config: config})
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatal(err)
	}
	parityJSON(t, p, string(raw))
}

func TestCapabilityParityParamNullMembers(t *testing.T) {
	for _, tc := range []struct {
		value any
		valid string
	}{
		{forward.EnvironmentConfigParams{}, `{"type":"cloud","packages":{"apt":["curl"]},"setup_script":"setup","startup_script":"startup","storage":{"buckets":[]}}`},
		{forward.EnvironmentConfigPackagesParam{}, `{"apt":["curl"],"npm":[],"pip":[],"cargo":[],"gem":[],"go":[]}`},
		{forward.EnvironmentStorageParam{}, `{"buckets":[{"name":"data","provider":"aliyun-oss","region":"cn-hangzhou","bucket":"example","credential":{"type":"ram_role","role_arn":"test-role"}}]}`},
		{forward.EnvironmentStorageBucketParam{}, `{"name":"data","provider":"aliyun-oss","region":"cn-hangzhou","bucket":"example","prefix":"team/","access":"read_only","accelerator_endpoint":"example.test","credential":{"type":"access_key","access_key_id":"test-id","access_key_secret":"test-secret"}}`},
		{forward.EnvironmentStorageCredentialParam{}, `{"type":"access_key","role_arn":"test-role","access_key_id":"test-id","access_key_secret":"test-secret"}`},
		{forward.TemplateA2AConfigParam{}, `{"enabled":true,"description":"Public","skills":[{"id":"summary","name":"Summary","description":"Public","tags":[],"examples":[]}]}`},
		{forward.TemplateA2ASkillParam{}, `{"id":"summary","name":"Summary","description":"Public","tags":[],"examples":[]}`},
	} {
		typ := reflect.TypeOf(tc.value)
		t.Run(typ.Name(), func(t *testing.T) {
			cases := map[string]string{"values": tc.valid, "empty object": `{}`, "whole null": `null`}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal([]byte(tc.valid), &fields); err != nil {
				t.Fatal(err)
			}
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				name := strings.Split(field.Tag.Get("json"), ",")[0]
				original, present := fields[name]
				if !present {
					continue
				}
				fields[name] = json.RawMessage(`null`)
				raw, err := json.Marshal(fields)
				if err != nil {
					t.Fatal(err)
				}
				cases[name+" null"] = string(raw)
				if field.Type.Kind() == reflect.Slice {
					fields[name] = json.RawMessage(`[null]`)
					raw, err = json.Marshal(fields)
					if err != nil {
						t.Fatal(err)
					}
					cases[name+" null item"] = string(raw)
					fields[name] = json.RawMessage(`[]`)
					raw, err = json.Marshal(fields)
					if err != nil {
						t.Fatal(err)
					}
					cases[name+" empty array"] = string(raw)
				}
				fields[name] = original
			}
			for name, raw := range cases {
				t.Run(name, func(t *testing.T) {
					for _, decoder := range []struct {
						name string
						call func([]byte, any) error
					}{
						{"json", json.Unmarshal},
						{"apijson", apijson.Unmarshal},
						{"APIJSON", func(raw []byte, dst any) error {
							return dst.(interface{ UnmarshalAPIJSON([]byte) error }).UnmarshalAPIJSON(raw)
						}},
					} {
						t.Run(decoder.name, func(t *testing.T) {
							decoded := reflect.New(typ)
							if err := decoder.call([]byte(raw), decoded.Interface()); err != nil {
								t.Fatal(err)
							}
							parityJSON(t, decoded.Interface(), raw)
							if raw == `{}` && param.IsOmitted(decoded.Elem().Interface()) {
								t.Fatal("authored empty object became omission")
							}
							if raw == `null` && !param.IsNull(decoded.Elem().Interface()) {
								t.Fatal("lost whole-object null sentinel")
							}
						})
					}
				})
			}
		})
	}
}

func TestCapabilityParityMixedNullExtras(t *testing.T) {
	extras := map[string]any{"future": 7, "future_null": nil}
	config := forward.EnvironmentConfigParams{}
	config.SetExtraFields(extras)
	if err := json.Unmarshal([]byte(`{"type":"cloud","storage":null}`), &config); err != nil {
		t.Fatal(err)
	}
	parityJSON(t, config, `{"type":"cloud","storage":null,"future":7,"future_null":null}`)
	if _, changed := extras["storage"]; changed {
		t.Fatal("null restoration mutated the caller's extra-fields map")
	}
	config.Type = forward.String("self_hosted")
	parityJSON(t, config, `{"type":"self_hosted","storage":null,"future":7,"future_null":null}`)
	if err := json.Unmarshal([]byte(`{"type":"cloud","storage":{"buckets":[]}}`), &config); err != nil {
		t.Fatal(err)
	}
	parityJSON(t, config, `{"type":"cloud","storage":{"buckets":[]},"future":7,"future_null":null}`)

	var a2a forward.TemplateA2AConfigParam
	if err := json.Unmarshal([]byte(`{"enabled":true,"skills":null}`), &a2a); err != nil {
		t.Fatal(err)
	}
	a2a.Enabled = forward.Bool(false)
	parityJSON(t, a2a, `{"enabled":false,"skills":null}`)
	if err := json.Unmarshal([]byte(`{"enabled":true,"skills":[]}`), &a2a); err != nil {
		t.Fatal(err)
	}
	parityJSON(t, a2a, `{"enabled":true,"skills":[]}`)
	var packages forward.EnvironmentConfigPackagesParam
	if err := json.Unmarshal([]byte(`{"apt":["curl",null]}`), &packages); err != nil {
		t.Fatal(err)
	}
	parityJSON(t, packages, `{"apt":["curl",null]}`)
	if err := json.Unmarshal([]byte(`{"apt":[]}`), &packages); err != nil {
		t.Fatal(err)
	}
	parityJSON(t, packages, `{"apt":[]}`)
	for _, tc := range []struct{ raw, want string }{
		{`{"enabled":true,"skills":null,"future":7}`, `{"enabled":true,"skills":null,"future":7}`},
		{`{"skills":null,"future":7}`, `{"skills":null,"future":7}`},
		{`{"future":7}`, `{"future":7}`},
	} {
		var decoded forward.TemplateA2AConfigParam
		if err := json.Unmarshal([]byte(tc.raw), &decoded); err != nil {
			t.Fatal(err)
		}
		parityJSON(t, decoded, tc.want)
	}
	var decoded forward.EnvironmentConfigParams
	if err := json.Unmarshal([]byte(`{"type":"cloud","storage":null,"future":7}`), &decoded); err != nil {
		t.Fatal(err)
	}
	parityJSON(t, decoded, `{"type":"cloud","storage":null,"future":7}`)
	bucket := forward.EnvironmentStorageBucketParam{Credential: forward.EnvironmentStorageCredentialParam{Type: "access_key", AccessKeySecret: forward.String("old-test-secret")}}
	raw := `{"name":"data","provider":"aliyun-oss","region":"cn-hangzhou","bucket":"example","credential":null}`
	if err := json.Unmarshal([]byte(raw), &bucket); err != nil {
		t.Fatal(err)
	}
	parityJSON(t, bucket, raw)
}

func TestCapabilityParityNestedParamNulls(t *testing.T) {
	for _, tc := range []struct {
		value any
		raws  []string
	}{
		{forward.EnvironmentNewParams{}, []string{`{"name":"env","config":null}`, `{"name":"env","config":{}}`, `{"name":"env","config":{"type":"cloud","storage":null}}`}},
		{forward.EnvironmentUpdateParams{}, []string{`{"config":null}`, `{"config":{}}`, `{"config":{"type":"cloud","packages":{"apt":null,"npm":[]}}}`}},
		{forward.TemplateNewParams{}, []string{`{"name":"tmpl","model":"ultimate","environment_id":"env","a2a":null}`, `{"name":"tmpl","model":"ultimate","environment_id":"env","a2a":{}}`, `{"name":"tmpl","model":"ultimate","environment_id":"env","a2a":{"enabled":true,"skills":null}}`}},
		{forward.TemplateUpdateParams{}, []string{`{"a2a":null}`, `{"a2a":{}}`, `{"a2a":{"enabled":true,"skills":null}}`}},
	} {
		for _, raw := range tc.raws {
			t.Run(reflect.TypeOf(tc.value).Name()+raw, func(t *testing.T) {
				for _, decode := range []func([]byte, any) error{json.Unmarshal, apijson.Unmarshal} {
					decoded := reflect.New(reflect.TypeOf(tc.value))
					if err := decode([]byte(raw), decoded.Interface()); err != nil {
						t.Fatal(err)
					}
					parityJSON(t, decoded.Interface(), raw)
				}
			})
		}
	}
}

func TestCapabilityParityConfigUnknownFields(t *testing.T) {
	for _, tc := range []struct {
		value any
		raw   string
	}{
		{&forward.EnvironmentConfigParams{}, `{"type":"cloud","networking":{"egress":"restricted"},"future":null}`},
		{&forward.EnvironmentConfigPackagesParam{}, `{"apt":[],"future":{"version":9007199254740993}}`},
		{&forward.EnvironmentStorageParam{}, `{"buckets":[],"future":true}`},
		{&forward.EnvironmentStorageBucketParam{}, `{"name":"data","provider":"aliyun-oss","region":"cn-hangzhou","bucket":"example","future":7}`},
		{&forward.EnvironmentStorageCredentialParam{}, `{"type":"ram_role","role_arn":"test-role","future":7}`},
		{&forward.TemplateA2AConfigParam{}, `{"enabled":true,"skills":[],"future":7}`},
		{&forward.TemplateA2ASkillParam{}, `{"id":"summary","name":"Summary","description":"Public","future":7}`},
	} {
		t.Run(reflect.TypeOf(tc.value).Elem().Name(), func(t *testing.T) {
			if err := json.Unmarshal([]byte(tc.raw), tc.value); err != nil {
				t.Fatal(err)
			}
			parityJSON(t, tc.value, tc.raw)
			if strings.Contains(tc.raw, "9007199254740993") {
				raw, err := json.Marshal(tc.value)
				if err != nil || !strings.Contains(string(raw), "9007199254740993") {
					t.Fatalf("unknown field lost numeric precision: %s, %v", raw, err)
				}
			}
		})
	}
	var config forward.EnvironmentConfigParams
	const input = `{"type":"cloud","networking":{"egress":"restricted"},"storage":{"buckets":[]}}`
	if err := json.Unmarshal([]byte(input), &config); err != nil {
		t.Fatal(err)
	}
	client := testClient(func(r *http.Request) (*http.Response, error) {
		var body any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		parityJSON(t, body, `{"config":{"type":"cloud","networking":{"egress":"restricted"},"storage":{"buckets":[]}}}`)
		return reply(r, 200, `{}`), nil
	})
	if _, err := client.Environments.Update(context.Background(), "target", forward.EnvironmentUpdateParams{Config: config}); err != nil {
		t.Fatal(err)
	}
}

func TestCapabilityParityInvalidMembersRemainVisible(t *testing.T) {
	for _, raw := range []string{
		`{"type":"cloud","storage":[]}`,
		`{"type":"cloud","packages":false}`,
		`{"type":"cloud","storage":{"buckets":[[]]}}`,
	} {
		var config forward.EnvironmentConfigParams
		if err := json.Unmarshal([]byte(raw), &config); err != nil {
			t.Fatal(err)
		}
		parityJSON(t, config, raw)
	}
}

func TestCapabilityParityTypedEditsReplaceDecodedRaw(t *testing.T) {
	var a2a forward.TemplateA2AConfigParam
	if err := json.Unmarshal([]byte(`{"description":null,"skills":null}`), &a2a); err != nil {
		t.Fatal(err)
	}
	a2a.Description = forward.String("")
	a2a.Skills = []forward.TemplateA2ASkillParam{}
	parityJSON(t, a2a, `{"description":"","skills":[]}`)

	var config forward.EnvironmentConfigParams
	if err := json.Unmarshal([]byte(`{"type":"cloud","storage":[]}`), &config); err != nil {
		t.Fatal(err)
	}
	config.Storage = forward.EnvironmentStorageParam{Buckets: []forward.EnvironmentStorageBucketParam{}}
	parityJSON(t, config, `{"type":"cloud","storage":{"buckets":[]}}`)
	config.Storage = param.Override[forward.EnvironmentStorageParam](json.RawMessage(`{}`))
	parityJSON(t, config, `{"type":"cloud","storage":{}}`)

	var credential forward.EnvironmentStorageCredentialParam
	if err := json.Unmarshal([]byte(`{"type":null}`), &credential); err != nil {
		t.Fatal(err)
	}
	credential.Type = "ram_role"
	credential.RoleARN = forward.String("test-role")
	parityJSON(t, credential, `{"type":"ram_role","role_arn":"test-role"}`)

	a2a.SetExtraFields(map[string]any{"description": nil})
	a2a.Description = forward.String("manual overrides remain authoritative")
	parityJSON(t, a2a, `{"description":null,"skills":[]}`)
	param.SetJSON([]byte(`{"enabled":false}`), &a2a)
	parityJSON(t, a2a, `{"enabled":false}`)
}

func TestCapabilityParityEditableEmptyObject(t *testing.T) {
	var config forward.EnvironmentConfigParams
	if err := json.Unmarshal([]byte(`{}`), &config); err != nil {
		t.Fatal(err)
	}
	config.Type = forward.String("cloud")
	config.Storage = forward.EnvironmentStorageParam{Buckets: []forward.EnvironmentStorageBucketParam{}}
	parityJSON(t, config, `{"type":"cloud","storage":{"buckets":[]}}`)

	var a2a forward.TemplateA2AConfigParam
	if err := json.Unmarshal([]byte(`{}`), &a2a); err != nil {
		t.Fatal(err)
	}
	a2a.Enabled = forward.Bool(true)
	a2a.Skills = []forward.TemplateA2ASkillParam{}
	parityJSON(t, a2a, `{"enabled":true,"skills":[]}`)

	var credential forward.EnvironmentStorageCredentialParam
	if err := json.Unmarshal([]byte(`{}`), &credential); err != nil {
		t.Fatal(err)
	}
	credential.Type = "ram_role"
	credential.RoleARN = forward.String("test-role")
	parityJSON(t, credential, `{"type":"ram_role","role_arn":"test-role"}`)

	var skill forward.TemplateA2ASkillParam
	if err := json.Unmarshal([]byte(`{}`), &skill); err != nil {
		t.Fatal(err)
	}
	skill.Name = "Summary"
	parityJSON(t, skill, `{"name":"Summary"}`)

	client := testClient(func(r *http.Request) (*http.Response, error) {
		var body any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		parityJSON(t, body, `{"a2a":{"enabled":true,"skills":[]}}`)
		return reply(r, 200, `{}`), nil
	})
	if _, err := client.Templates.Update(context.Background(), "target", forward.TemplateUpdateParams{A2A: a2a}); err != nil {
		t.Fatal(err)
	}
}

func TestCapabilityParitySessionRosterReuse(t *testing.T) {
	var config forward.SessionUpdateParamsConfigParam
	if err := json.Unmarshal([]byte(`{"multiagent":null}`), &config); err != nil {
		t.Fatal(err)
	}
	if !param.IsNull(config.Multiagent) {
		t.Fatal("missing initial null")
	}
	const raw = `{"multiagent":{"type":"coordinator","agents":["tmpl_child"]}}`
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		t.Fatal(err)
	}
	parityJSON(t, config, raw)
	client := testClient(func(r *http.Request) (*http.Response, error) {
		var body any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		parityJSON(t, body, `{"config":{"multiagent":{"type":"coordinator","agents":["tmpl_child"]}}}`)
		return reply(r, 200, `{}`), nil
	})
	if _, err := client.Sessions.Update(context.Background(), "target", forward.SessionUpdateParams{Config: config}); err != nil {
		t.Fatal(err)
	}
}

func TestCapabilityParityMixedNullAPIContracts(t *testing.T) {
	for _, surface := range []struct {
		name, resource, field string
		cases                 []struct {
			name, value string
			status      int
		}
	}{
		{"Environment", "environments", "config", []struct {
			name, value string
			status      int
		}{
			{"omitted", "", 200},
			{"whole null", `null`, 400},
			{"storage null", `{"type":"cloud","storage":null}`, 400},
			{"storage array", `{"type":"cloud","storage":[]}`, 400},
			{"packages null", `{"type":"cloud","packages":null}`, 400},
			{"apt null", `{"type":"cloud","packages":{"apt":null,"npm":[]}}`, 400},
			{"apt null item", `{"type":"cloud","packages":{"apt":["curl",null]}}`, 400},
			{"buckets null", `{"type":"cloud","storage":{"buckets":null}}`, 400},
			{"credential null", `{"type":"cloud","storage":{"buckets":[{"name":"data","provider":"aliyun-oss","region":"cn-hangzhou","bucket":"example","credential":null}]}}`, 400},
			{"empty buckets", `{"type":"cloud","storage":{"buckets":[]}}`, 200},
			{"empty packages", `{"type":"cloud","packages":{"apt":[],"npm":[]}}`, 200},
		}},
		{"Template", "templates", "a2a", []struct {
			name, value string
			status      int
		}{
			{"omitted", "", 200},
			{"whole null", `null`, 200},
			{"empty object", `{}`, 200},
			{"skills null", `{"enabled":true,"skills":null}`, 400},
			{"enabled null", `{"enabled":null,"skills":[]}`, 400},
			{"description null", `{"enabled":true,"description":null}`, 400},
			{"id null", `{"enabled":true,"skills":[{"id":null,"name":"Summary","description":"Public"}]}`, 400},
			{"tags null", `{"enabled":true,"skills":[{"id":"summary","name":"Summary","description":"Public","tags":null}]}`, 400},
			{"examples null item", `{"enabled":true,"skills":[{"id":"summary","name":"Summary","description":"Public","examples":[null]}]}`, 400},
			{"empty skills", `{"enabled":true,"skills":[]}`, 200},
		}},
	} {
		for _, operation := range []string{"create", "update"} {
			for _, tc := range surface.cases {
				t.Run(surface.name+"/"+operation+"/"+tc.name, func(t *testing.T) {
					body := map[string]json.RawMessage{}
					if operation == "create" {
						body["name"] = json.RawMessage(`"roundtrip"`)
						if surface.name == "Template" {
							body["model"] = json.RawMessage(`"ultimate"`)
							body["environment_id"] = json.RawMessage(`"env"`)
						}
					}
					if tc.value != "" {
						body[surface.field] = json.RawMessage(tc.value)
					}
					raw, err := json.Marshal(body)
					if err != nil {
						t.Fatal(err)
					}
					path := "/api/v1/forward/" + surface.resource
					if operation == "update" {
						path += "/target"
					}
					calls := 0
					client := testClient(func(r *http.Request) (*http.Response, error) {
						calls++
						if r.Method != http.MethodPost || r.URL.Path != path {
							t.Fatal(r.Method, r.URL)
						}
						var sent any
						if err := json.NewDecoder(r.Body).Decode(&sent); err != nil {
							t.Fatal(err)
						}
						parityJSON(t, sent, string(raw))
						if tc.status == 400 {
							return reply(r, 400, `{"error":{"type":"invalid_request_error","message":"invalid null member"}}`), nil
						}
						return reply(r, 200, `{}`), nil
					})
					switch surface.name + operation {
					case "Environmentcreate":
						var p forward.EnvironmentNewParams
						if err = json.Unmarshal(raw, &p); err == nil {
							parityJSON(t, p, string(raw))
							_, err = client.Environments.New(context.Background(), p)
						}
					case "Environmentupdate":
						var p forward.EnvironmentUpdateParams
						if err = json.Unmarshal(raw, &p); err == nil {
							parityJSON(t, p, string(raw))
							_, err = client.Environments.Update(context.Background(), "target", p)
						}
					case "Templatecreate":
						var p forward.TemplateNewParams
						if err = json.Unmarshal(raw, &p); err == nil {
							parityJSON(t, p, string(raw))
							_, err = client.Templates.New(context.Background(), p)
						}
					case "Templateupdate":
						var p forward.TemplateUpdateParams
						if err = json.Unmarshal(raw, &p); err == nil {
							parityJSON(t, p, string(raw))
							_, err = client.Templates.Update(context.Background(), "target", p)
						}
					}
					if tc.status == 400 {
						var apiErr *forward.Error
						if !errors.As(err, &apiErr) || apiErr.StatusCode != 400 || apiErr.Message != "invalid null member" {
							t.Fatal("lost service validation error", err)
						}
					} else if err != nil {
						t.Fatal(err)
					}
					if calls != 1 {
						t.Fatal("validation must remain server-owned", calls, err)
					}
				})
			}
		}
	}
}
