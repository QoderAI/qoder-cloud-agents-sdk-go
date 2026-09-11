package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestDeploymentNew(t *testing.T) {
	client := contractClient(t, "DeploymentService", "New")
	response, err := client.Deployments.New(context.Background(), managed.DeploymentNewParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Deployment.New")
}

func TestDeploymentGet(t *testing.T) {
	client := contractClient(t, "DeploymentService", "Get")
	response, err := client.Deployments.Get(context.Background(), "segment /?%#", managed.DeploymentGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Deployment.Get")
}

func TestDeploymentUpdate(t *testing.T) {
	client := contractClient(t, "DeploymentService", "Update")
	response, err := client.Deployments.Update(context.Background(), "segment /?%#", managed.DeploymentUpdateParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Deployment.Update")
}

func TestDeploymentList(t *testing.T) {
	client := contractClient(t, "DeploymentService", "List")
	response, err := client.Deployments.List(context.Background(), managed.DeploymentListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Deployment.List")
}

func TestDeploymentArchive(t *testing.T) {
	client := contractClient(t, "DeploymentService", "Archive")
	response, err := client.Deployments.Archive(context.Background(), "segment /?%#", managed.DeploymentArchiveParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Deployment.Archive")
}

func TestDeploymentPause(t *testing.T) {
	client := contractClient(t, "DeploymentService", "Pause")
	response, err := client.Deployments.Pause(context.Background(), "segment /?%#", managed.DeploymentPauseParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Deployment.Pause")
}

func TestDeploymentRun(t *testing.T) {
	client := contractClient(t, "DeploymentService", "Run")
	response, err := client.Deployments.Run(context.Background(), "segment /?%#", managed.DeploymentRunParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Deployment.Run")
}

func TestDeploymentUnpause(t *testing.T) {
	client := contractClient(t, "DeploymentService", "Unpause")
	response, err := client.Deployments.Unpause(context.Background(), "segment /?%#", managed.DeploymentUnpauseParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Deployment.Unpause")
}
