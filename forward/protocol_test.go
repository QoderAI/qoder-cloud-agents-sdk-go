package forward_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestOptionalNullAndUnionParameters(t *testing.T) {
	params := forward.TemplateUpdateParams{
		Name: forward.String(""), System: param.Null[string](), Tools: []forward.ToolParam{},
		Model:                forward.ModelConfigUnionParam{OfConfig: &forward.ModelConfigParam{ID: "ultimate", Effort: forward.String("high"), ContextWindow: forward.Int(400000)}},
		Multiagent:           param.NullStruct[forward.MultiagentConfigParam](),
		EnvironmentVariables: forward.EnvironmentVariablesUnionParam{OfMap: map[string]any{"CUSTOM_VAR": "value", "REMOVE": nil}},
		Metadata:             param.NullMap[map[string]any](),
	}
	params.SetExtraFields(map[string]any{"future.flag": false})
	obj := jsonObject(t, params)
	if obj["name"] != "" || obj["system"] != nil || obj["multiagent"] != nil || obj["metadata"] != nil || obj["future.flag"] != false {
		t.Fatalf("null/zero/extra fields: %#v", obj)
	}
	if _, ok := obj["description"]; ok {
		t.Fatal("omitted description was sent")
	}
	if tools, ok := obj["tools"].([]any); !ok || len(tools) != 0 {
		t.Fatalf("explicit empty slice: %#v", obj["tools"])
	}
	if model := obj["model"].(map[string]any); model["id"] != "ultimate" || model["context_window"] != float64(400000) {
		t.Fatal(model)
	}
	for _, input := range []string{`"ultimate"`, `{"id":"ultimate","effort":"high","future":true}`} {
		var model forward.ModelConfigUnionParam
		if err := json.Unmarshal([]byte(input), &model); err != nil {
			t.Fatal(err)
		}
		encoded, err := json.Marshal(model)
		if err != nil {
			t.Fatal(err)
		}
		var got, want any
		json.Unmarshal(encoded, &got)
		json.Unmarshal([]byte(input), &want)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("model union dropped fields: %s -> %s", input, encoded)
		}
	}
	_, err := json.Marshal(forward.ModelConfigUnionParam{OfString: forward.String("one"), OfConfig: &forward.ModelConfigParam{ID: "two"}})
	if err == nil {
		t.Fatal("ambiguous union accepted")
	}
}

func TestIdentityOverridesPreserveDynamicKeys(t *testing.T) {
	params := forward.IdentityConfigUpsertParams{IdentityConfig: forward.IdentityConfigSpecParam{
		Tools:                map[string]forward.ToolOverrideParam{"Read": {Enabled: forward.Bool(false)}},
		Skills:               map[string]forward.SkillOverrideParam{"skill_one": {Enabled: forward.Bool(false)}, "skill_inherit": param.NullStruct[forward.SkillOverrideParam]()},
		MCPServers:           map[string]forward.MCPServerOverrideParam{"crm": {Enabled: forward.Bool(true), URL: forward.String("https://crm.test/mcp")}},
		GitHubRepositories:   map[string]forward.GitHubRepositoryParam{"custom_binding": {AuthorizationToken: forward.String("test_token"), Enabled: param.Null[bool]()}},
		EnvironmentVariables: map[string]forward.EnvironmentVariableOverrideParam{"FOO": {Op: "set", Value: forward.String("value")}, "BAR": {Op: "unset"}},
	}}
	obj := jsonObject(t, params)["identity_config"].(map[string]any)
	tools := obj["tools"].(map[string]any)["Read"].(map[string]any)
	if len(tools) != 1 || tools["enabled"] != false {
		t.Fatalf("keyed override injected a name: %#v", tools)
	}
	skills := obj["skills"].(map[string]any)
	if skills["skill_inherit"] != nil || len(skills["skill_one"].(map[string]any)) != 1 {
		t.Fatal(skills)
	}
	if obj["environment_variables"].(map[string]any)["BAR"].(map[string]any)["op"] != "unset" {
		t.Fatal(obj)
	}
	var config forward.IdentityConfig
	if err := json.Unmarshal([]byte(`{"identity_config":{"skills":{"another":{"enabled":false}},"environment_variables":{"RUNTIME_KEY":{"op":"set","value":"dynamic"}}},"metadata":{"a.b":{"nested":1}}}`), &config); err != nil {
		t.Fatal(err)
	}
	checkDecoded(t, config)
	if config.IdentityConfig.EnvironmentVariables["RUNTIME_KEY"].Value != "dynamic" || config.IdentityConfig.Skills["another"].Enabled {
		t.Fatalf("dynamic map decode: %+v", config.IdentityConfig)
	}
}

func TestResourceResponseFieldPresence(t *testing.T) {
	var template forward.Template
	input := `{"id":"tmpl_one","model":"ultimate","environment_variables":{"KEY_NOT_IN_EXAMPLE":"value"},"metadata":{"extra":true},"description":null,"new_flag":true}`
	if err := json.Unmarshal([]byte(input), &template); err != nil {
		t.Fatal(err)
	}
	if template.ID != "tmpl_one" || template.Model.ID != "ultimate" || template.Model.RawJSON() != `"ultimate"` {
		t.Fatalf("template: %+v", template)
	}
	if template.EnvironmentVariables["KEY_NOT_IN_EXAMPLE"] != "value" {
		t.Fatal(template.EnvironmentVariables)
	}
	if template.JSON.Description.Raw() != "null" || template.JSON.Name.Raw() != "" || template.JSON.ExtraFields["new_flag"].Raw() != "true" {
		t.Fatal("presence/null/unknown fields lost")
	}
	if template.RawJSON() != input {
		t.Fatal("raw response changed")
	}
	checkDecoded(t, template)
}

func TestEnvironmentPackagesAndInitialResourceMount(t *testing.T) {
	var environment forward.Environment
	if err := json.Unmarshal([]byte(`{"config":{"type":"sandbox","packages":{"apt":["git","curl"],"npm":["pnpm@9"],"pip":["PyYAML==6.0.1"],"cargo":["ripgrep"],"gem":["bundler"],"go":["example.com/tool@latest"]}}}`), &environment); err != nil {
		t.Fatal(err)
	}
	checkDecoded(t, environment)
	if !reflect.DeepEqual(environment.Config.Packages.Apt, []string{"git", "curl"}) || environment.Config.Packages.Pip[0] != "PyYAML==6.0.1" {
		t.Fatal(environment.Config.Packages)
	}
	var params forward.SessionNewParams
	if err := json.Unmarshal([]byte(`{"identity_id":"idn_one","template_id":"tmpl_one","resources":[{"type":"file","file_id":"file_one","mount_path":"/data/input.py"}]}`), &params); err != nil {
		t.Fatal(err)
	}
	if len(params.Resources) != 1 || params.Resources[0].MountPath.Value != "/data/input.py" {
		t.Fatalf("mount path lost: %#v", params.Resources)
	}
	if obj := jsonObject(t, params); obj["resources"].([]any)[0].(map[string]any)["mount_path"] != "/data/input.py" {
		t.Fatal(obj)
	}
}

func TestSkillMultipartAndTypedResponse(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		reader, err := r.MultipartReader()
		if err != nil {
			t.Fatal(err)
		}
		files := map[string]string{}
		for {
			p, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			b, _ := io.ReadAll(p)
			if p.FormName() == "files" {
				files[p.Header.Get("Content-Disposition")] = string(b)
			}
			if p.FormName() == "metadata" && string(b) != `{"team":"sdk"}` {
				t.Errorf("metadata: %s", b)
			}
		}
		if len(files) != 2 {
			t.Fatalf("file parts: %v", files)
		}
		found := false
		for name := range files {
			found = found || strings.Contains(name, "skill/scripts/run.sh")
		}
		if !found {
			t.Fatal("relative filename was lost")
		}
		return reply(r, 201, `{"id":"skill_one","display_title":"skill","latest_version":"123"}`), nil
	})
	skill, err := client.Skills.New(context.Background(), forward.SkillNewParams{Files: []io.Reader{
		convention.UploadFile{Reader: strings.NewReader("# Skill"), Name: "skill/SKILL.md"},
		convention.UploadFile{Reader: strings.NewReader("echo test"), Name: "skill/scripts/run.sh"},
	}, Metadata: map[string]any{"team": "sdk"}})
	if err != nil || skill.ID != "skill_one" || skill.LatestVersion != "123" || calls != 1 {
		t.Fatalf("skill: %+v %v", skill, err)
	}
}

func TestDownloadDoesNotForwardAPIHeadersOrMiddleware(t *testing.T) {
	calls, middlewareCalls := 0, 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Host == "qoder.test" {
			if r.Header.Get("X-Private") != "api-only" {
				t.Fatal("API middleware missing")
			}
			return reply(r, 200, `{"url":"https://storage.test/object?signature=signed"}`), nil
		}
		if r.Header.Get("Authorization") != "" || r.Header.Get("X-Private") != "" || r.URL.RawQuery != "signature=signed" {
			t.Fatalf("storage request: %s %v", r.URL, r.Header)
		}
		return reply(r, 200, "content"), nil
	}, option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		middlewareCalls++
		r.Header.Set("X-Private", "api-only")
		return next(r)
	}))
	response, err := client.Files.Download(context.Background(), "file_one")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil || string(b) != "content" || calls != 2 || middlewareCalls != 1 {
		t.Fatalf("download: %q %v calls=%d middleware=%d", b, err, calls, middlewareCalls)
	}
}

func TestSSEPreservesDeltasAndUnknownEvents(t *testing.T) {
	wire := `id: evt_one
event: event_start
data: {"id":"evt_one","type":"event_start","event":{"id":"evt_one","type":"agent.message"}}

id: evt_one
event: event_delta
data: {"id":"evt_one","type":"event_delta","event_id":"evt_one","delta":{"type":"content_delta","content":{"type":"text","text":"hi"}}}

id: evt_two
event: future.event
data: {"id":"evt_two","type":"future.event","new_field":{"value":1}}

`
	client := testClient(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Query()["event_deltas[]"]; !reflect.DeepEqual(got, []string{"agent.message", "agent.thinking"}) {
			t.Fatal(got)
		}
		if r.Header.Get("Last-Event-ID") != "evt_previous" {
			t.Fatal("resume header missing")
		}
		res := reply(r, 200, wire)
		res.Header.Set("Content-Type", "text/event-stream")
		return res, nil
	})
	stream := client.Sessions.Events.StreamEvents(context.Background(), "sess_one", forward.SessionEventStreamParams{EventDeltas: []string{"agent.message", "agent.thinking"}, LastEventID: forward.String("evt_previous")})
	defer stream.Close()
	var types []string
	for stream.Next() {
		event := stream.Current()
		types = append(types, event.Type)
		if event.Type == "event_delta" && !strings.Contains(string(event.Delta), `"hi"`) {
			t.Fatal("delta payload missing")
		}
	}
	if stream.Err() != nil || !reflect.DeepEqual(types, []string{"event_start", "event_delta", "future.event"}) || stream.LastEventID() != "evt_two" {
		t.Fatalf("stream: %v %v", types, stream.Err())
	}
	if stream.Current().JSON.ExtraFields["new_field"].Raw() != `{"value":1}` {
		t.Fatal("unknown event payload lost")
	}
}

func TestStreamingCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"id\":\"evt_one\",\"type\":\"agent.message\"}\n\n")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := forward.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithMaxRetries(0))
	stream := client.Sessions.Events.StreamEvents(ctx, "sess_one", forward.SessionEventStreamParams{})
	defer stream.Close()
	if !stream.Next() {
		t.Fatal(stream.Err())
	}
	cancel()
	done := make(chan struct{})
	go func() { defer close(done); stream.Next() }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("stream did not cancel")
	}
	if !errors.Is(stream.Err(), context.Canceled) {
		t.Fatalf("cancellation error: %v", stream.Err())
	}
}

func TestForwardRetrySafety(t *testing.T) {
	for _, tc := range []struct {
		name         string
		status, want int
		write, key   bool
	}{
		{"get_503", 503, 3, false, false}, {"write_503", 503, 1, true, false}, {"idempotent_write_503", 503, 3, true, true}, {"write_429", 429, 3, true, false}, {"write_409", 409, 1, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				res := reply(r, tc.status, `{"error":{"type":"api_error","message":"temporary"}}`)
				res.Header.Set("Retry-After", "0")
				return res, nil
			}, option.WithMaxRetries(2))
			var err error
			if tc.write {
				p := forward.TemplateNewParams{Name: "test", EnvironmentID: "env", Model: forward.ModelConfigUnionParam{OfString: forward.String("model")}}
				if tc.key {
					p.IdempotencyKey = forward.String("key")
				}
				_, err = client.Templates.New(context.Background(), p)
			} else {
				_, err = client.Models.List(context.Background())
			}
			if err == nil || calls != tc.want {
				t.Fatalf("calls=%d want=%d err=%v", calls, tc.want, err)
			}
		})
	}
}
