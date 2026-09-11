package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestSkillVersionNew(t *testing.T) {
	client := contractClient(t, "SkillVersionService", "New")
	response, err := client.Skills.Versions.New(context.Background(), "segment /?%#", managed.SkillVersionNewParams{Files: []io.Reader{strings.NewReader("file")}})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SkillVersion.New")
}

func TestSkillVersionGet(t *testing.T) {
	client := contractClient(t, "SkillVersionService", "Get")
	response, err := client.Skills.Versions.Get(context.Background(), "segment /?%#", managed.SkillVersionGetParams{SkillID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SkillVersion.Get")
}

func TestSkillVersionList(t *testing.T) {
	client := contractClient(t, "SkillVersionService", "List")
	response, err := client.Skills.Versions.List(context.Background(), "segment /?%#", managed.SkillVersionListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SkillVersion.List")
}

func TestSkillVersionDelete(t *testing.T) {
	client := contractClient(t, "SkillVersionService", "Delete")
	response, err := client.Skills.Versions.Delete(context.Background(), "segment /?%#", managed.SkillVersionDeleteParams{SkillID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SkillVersion.Delete")
}

func TestSkillVersionDownload(t *testing.T) {
	client := contractClient(t, "SkillVersionService", "Download")
	response, err := client.Skills.Versions.Download(context.Background(), "segment /?%#", managed.SkillVersionDownloadParams{SkillID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
}
