package managed_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func alignmentFixture(t *testing.T) map[string]json.RawMessage {
	t.Helper()
	data, err := os.ReadFile("testdata/contract-alignment-0.8.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture map[string]json.RawMessage
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	return fixture
}

func TestManagedContractAlignment080Resources(t *testing.T) {
	f := alignmentFixture(t)
	git := managed.ManagedAgentsGitRepositoryResourceParams{Type: "git_repository", URL: "https://user@git.example.test/team/repo.git", Password: managed.String("contract-placeholder"), MountPath: managed.String("/repos/demo"), Checkout: managed.ManagedAgentsGitRepositoryResourceParamsCheckoutUnion{OfBranch: &managed.ManagedAgentsBranchCheckoutParam{Type: "branch", Name: "main"}}}
	github := managed.ManagedAgentsGitHubRepositoryResourceParams{Type: "github_repository", URL: "https://github.com/QoderAI/forward-quickstart"}
	var expected []any
	if err := json.Unmarshal([]byte("["+string(f["git"])+","+string(f["github"])+"]"), &expected); err != nil {
		t.Fatal(err)
	}
	client := testClient(func(r *http.Request) (*http.Response, error) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]any
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(r.URL.Path, "environments") {
			if !reflect.DeepEqual(body["config"], map[string]any{"type": "self_hosted", "setup_script": "echo sdk-contract"}) {
				t.Fatal(body)
			}
		} else if !reflect.DeepEqual(body["resources"], expected) {
			t.Fatalf("resources=%s", data)
		}
		return reply(r, 200, `{"id":"contract_one"}`), nil
	})
	ctx := context.Background()
	_, err := client.Sessions.New(ctx, managed.SessionNewParams{Agent: managed.SessionNewParamsAgentUnion{OfString: managed.String("agent_one")}, EnvironmentID: "env_one", Resources: []managed.SessionNewParamsResourceUnion{{OfGitRepository: &git}, {OfGitHubRepository: &github}}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Deployments.New(ctx, managed.DeploymentNewParams{Resources: []managed.DeploymentNewParamsResourceUnion{{OfGitRepository: &git}, {OfGitHubRepository: &github}}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Deployments.Update(ctx, "dep_one", managed.DeploymentUpdateParams{Resources: []managed.DeploymentUpdateParamsResourceUnion{{OfGitRepository: &git}, {OfGitHubRepository: &github}}})
	if err != nil {
		t.Fatal(err)
	}
	config := managed.SelfHostedConfigParams{Type: "self_hosted", SetupScript: managed.String("echo sdk-contract")}
	_, err = client.Environments.New(ctx, managed.EnvironmentNewParams{Name: "sdk-contract", Config: managed.EnvironmentNewParamsConfigUnion{OfSelfHosted: &config}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Environments.Update(ctx, "env_one", managed.EnvironmentUpdateParams{Config: managed.EnvironmentUpdateParamsConfigUnion{OfSelfHosted: &config}})
	if err != nil {
		t.Fatal(err)
	}
	for _, dst := range []any{&managed.ManagedAgentsSessionResourceUnion{}, &managed.SessionResourceGetResponseUnion{}, &managed.SessionResourceUpdateResponseUnion{}, &managed.ManagedAgentsSessionResourceConfigUnion{}} {
		if err := json.Unmarshal([]byte(`{"type":"git_repository","url":"https://git.example.test/repo.git","mount_path":"/repos/demo","checkout":{"type":"branch","name":"main"}}`), dst); err != nil {
			t.Fatal(err)
		}
		switch v := dst.(type) {
		case *managed.ManagedAgentsSessionResourceUnion:
			if _, ok := v.AsAny().(managed.ManagedAgentsGitRepositoryResource); !ok {
				t.Fatalf("%T", v.AsAny())
			}
		case *managed.SessionResourceGetResponseUnion:
			if v.AsGitRepository().Checkout.AsBranch().Name != "main" {
				t.Fatal(v)
			}
		case *managed.SessionResourceUpdateResponseUnion:
			if v.AsGitRepository().MountPath != "/repos/demo" {
				t.Fatal(v)
			}
		case *managed.ManagedAgentsSessionResourceConfigUnion:
			if _, ok := v.AsAny().(managed.ManagedAgentsGitRepositoryResourceConfig); !ok {
				t.Fatalf("%T", v.AsAny())
			}
		}
	}
}

func TestManagedContractAlignment080ResponsesAndStreams(t *testing.T) {
	f := alignmentFixture(t)
	var memory managed.ManagedAgentsMemoryStore
	if err := json.Unmarshal(f["memory_store"], &memory); err != nil {
		t.Fatal(err)
	}
	if memory.Status != "active" || !memory.JSON.EntryCount.Valid() || memory.EntryCount != 0 || !memory.JSON.TotalSize.Valid() {
		t.Fatal(memory)
	}
	var dream managed.Dream
	if err := json.Unmarshal(f["dream"], &dream); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(dream.Outputs[0].FilesTouched, []string{"notes.md", "summary.md"}) {
		t.Fatal(dream)
	}
	var env managed.Environment
	if err := json.Unmarshal(f["environment"], &env); err != nil {
		t.Fatal(err)
	}
	if env.Config.SetupScript != "echo sdk-contract" || env.Config.AsSelfHosted().SetupScript != "echo sdk-contract" {
		t.Fatal(env)
	}
	var events []json.RawMessage
	if err := json.Unmarshal(f["events"], &events); err != nil {
		t.Fatal(err)
	}
	var wire strings.Builder
	for _, e := range events {
		var compact bytes.Buffer
		if err := json.Compact(&compact, e); err != nil {
			t.Fatal(err)
		}
		wire.WriteString("data: " + compact.String() + "\n\n")
	}
	client := testClient(func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "/stream") {
			res := reply(r, 200, wire.String())
			res.Header.Set("Content-Type", "text/event-stream")
			return res, nil
		}
		return reply(r, 200, `{"data":`+string(f["events"])+`,"has_more":false}`), nil
	})
	checkArtifact := func(v managed.ManagedAgentsAgentArtifactDeliveredEvent) {
		t.Helper()
		if v.FileID != "file_one" || v.Size != 0 || !v.JSON.Size.Valid() || v.OriginalFilename != "report.txt" || v.ContentType != "text/plain" {
			t.Fatal(v)
		}
	}
	page, err := client.Sessions.Events.List(context.Background(), "sess_one", managed.SessionEventListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkArtifact(page.Data[0].AsAgentArtifactDelivered())
	if _, ok := page.Data[0].AsAny().(managed.ManagedAgentsAgentArtifactDeliveredEvent); !ok {
		t.Fatal("history discriminator")
	}
	if page.Data[1].Error.QoderErrorCode != "TEST_DIAGNOSTIC" || page.Data[1].AsSessionError().Error.AsModelRequestFailedError().QoderErrorCode != "TEST_DIAGNOSTIC" {
		t.Fatal(page.Data[1])
	}
	stream := client.Sessions.Events.StreamEvents(context.Background(), "sess_one", managed.SessionEventStreamParams{})
	defer stream.Close()
	count := 0
	for stream.Next() {
		e := stream.Current()
		if count == 0 {
			checkArtifact(e.AsAgentArtifactDelivered())
		} else if e.Error.QoderErrorCode != "TEST_DIAGNOSTIC" || e.Error.RetryStatus.Type != "terminal" {
			t.Fatal(e)
		}
		count++
	}
	if stream.Err() != nil || count != 2 {
		t.Fatalf("stream count=%d err=%v", count, stream.Err())
	}
	thread := client.Sessions.Threads.Events.StreamEvents(context.Background(), "sthr_one", managed.SessionThreadEventStreamParams{SessionID: "sess_one"})
	defer thread.Close()
	count = 0
	for thread.Next() {
		e := thread.Current()
		if count == 0 {
			checkArtifact(e.AsAgentArtifactDelivered())
		} else if e.Error.QoderErrorCode != "TEST_DIAGNOSTIC" {
			t.Fatal(e)
		}
		count++
	}
	if thread.Err() != nil || count != 2 {
		t.Fatalf("thread count=%d err=%v", count, thread.Err())
	}
}

func TestManagedContractAlignment080CredentialRotation(t *testing.T) {
	client := testClient(func(r *http.Request) (*http.Response, error) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]any
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(body, map[string]any{"password": " token with spaces "}) {
			t.Fatal(body)
		}
		return reply(r, 200, `{"type":"git_repository","url":"https://user@git.example.test/repo.git"}`), nil
	})
	resource, err := client.Sessions.Resources.Update(context.Background(), "sesr_one", managed.SessionResourceUpdateParams{SessionID: "sess_one", Password: managed.String(" token with spaces ")})
	if err != nil {
		t.Fatal(err)
	}
	if resource.AsGitRepository().URL != "https://user@git.example.test/repo.git" {
		t.Fatal(resource)
	}
}
