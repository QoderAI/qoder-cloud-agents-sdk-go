package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestAgentVersionList(t *testing.T) {
	client := contractClient(t, "AgentVersionService", "List")
	response, err := client.Agents.Versions.List(context.Background(), "segment /?%#", managed.AgentVersionListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "AgentVersion.List")
}
