package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestMemoryStoreMemoryNew(t *testing.T) {
	client := contractClient(t, "MemoryStoreMemoryService", "New")
	response, err := client.MemoryStores.Memories.New(context.Background(), "segment /?%#", managed.MemoryStoreMemoryNewParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStoreMemory.New")
}

func TestMemoryStoreMemoryGet(t *testing.T) {
	client := contractClient(t, "MemoryStoreMemoryService", "Get")
	response, err := client.MemoryStores.Memories.Get(context.Background(), "segment /?%#", managed.MemoryStoreMemoryGetParams{MemoryStoreID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStoreMemory.Get")
}

func TestMemoryStoreMemoryUpdate(t *testing.T) {
	client := contractClient(t, "MemoryStoreMemoryService", "Update")
	response, err := client.MemoryStores.Memories.Update(context.Background(), "segment /?%#", managed.MemoryStoreMemoryUpdateParams{MemoryStoreID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStoreMemory.Update")
}

func TestMemoryStoreMemoryList(t *testing.T) {
	client := contractClient(t, "MemoryStoreMemoryService", "List")
	response, err := client.MemoryStores.Memories.List(context.Background(), "segment /?%#", managed.MemoryStoreMemoryListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStoreMemory.List")
}

func TestMemoryStoreMemoryDelete(t *testing.T) {
	client := contractClient(t, "MemoryStoreMemoryService", "Delete")
	response, err := client.MemoryStores.Memories.Delete(context.Background(), "segment /?%#", managed.MemoryStoreMemoryDeleteParams{MemoryStoreID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStoreMemory.Delete")
}

func TestQoderExtensionParameters(t *testing.T) {
	d := managed.DeploymentUpdateParams{EnvironmentVariables: param.Null[string]()}
	if v, ok := jsonObject(t, d)["environment_variables"]; !ok || v != nil {
		t.Fatal("deployment env clear lost")
	}
	s := managed.SessionNewParams{EnvironmentVariables: map[string]string{"ENV": "prod"}}
	if jsonObject(t, s)["environment_variables"].(map[string]any)["ENV"] != "prod" {
		t.Fatal("session creation env lost")
	}
	memory := managed.MemoryStoreMemoryUpdateParams{Content: managed.String("new"), Precondition: managed.ManagedAgentsPreconditionParam{Type: "content_sha256", ContentSha256: managed.String("expected-hash")}, Metadata: map[string]any{"team": "sdk"}}
	body := jsonObject(t, memory)
	if body["content_sha256"] != "expected-hash" {
		t.Fatal("optimistic concurrency precondition lost")
	}
	if _, ok := body["precondition"]; ok {
		t.Fatal("unconverted memory precondition sent")
	}
	stop := managed.EnvironmentWorkStopParams{SelfHostedWorkStopRequest: managed.SelfHostedWorkStopRequestParam{Force: managed.Bool(true)}}
	if jsonObject(t, stop)["force"] != true {
		t.Fatal("force stop lost")
	}
	add := managed.SessionResourceAddParams{ManagedAgentsFileResourceParams: managed.ManagedAgentsFileResourceParams{Type: "file", FileID: "file", MountPath: managed.String("/input")}}
	add.SetExtraFields(map[string]any{"mount_path": nil})
	body = jsonObject(t, add)
	if body["file_id"] != "file" || body["type"] != "file" {
		t.Fatal("resource payload nested")
	}
	if v, ok := body["mount_path"]; !ok || v != nil {
		t.Fatal("resource extra fields lost")
	}
	q, e := (managed.FileListParams{AfterID: managed.String("cursor"), Name: managed.String("report")}).URLQuery()
	if e != nil || q.Get("after_id") != "cursor" || q.Get("name") != "report" {
		t.Fatal(q, e)
	}
}
