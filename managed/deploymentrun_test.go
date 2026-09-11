package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestDeploymentRunGet(t *testing.T) {
	client := contractClient(t, "DeploymentRunService", "Get")
	response, err := client.DeploymentRuns.Get(context.Background(), "segment /?%#", managed.DeploymentRunGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "DeploymentRun.Get")
}

func TestDeploymentRunList(t *testing.T) {
	client := contractClient(t, "DeploymentRunService", "List")
	response, err := client.DeploymentRuns.List(context.Background(), managed.DeploymentRunListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "DeploymentRun.List")
}
