package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestModelList(t *testing.T) {
	client := contractClient(t, "ModelService", "List")
	response, err := client.Models.List(context.Background(), managed.ModelListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Model.List")
}
