package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestFileList(t *testing.T) {
	client := contractClient(t, "FileService", "List")
	response, err := client.Files.List(context.Background(), managed.FileListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "File.List")
}

func TestFileDelete(t *testing.T) {
	client := contractClient(t, "FileService", "Delete")
	response, err := client.Files.Delete(context.Background(), "segment /?%#", managed.FileDeleteParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "File.Delete")
}

func TestFileDownload(t *testing.T) {
	client := contractClient(t, "FileService", "Download")
	response, err := client.Files.Download(context.Background(), "segment /?%#", managed.FileDownloadParams{})
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
}

func TestFileGetMetadata(t *testing.T) {
	client := contractClient(t, "FileService", "GetMetadata")
	response, err := client.Files.GetMetadata(context.Background(), "segment /?%#", managed.FileGetMetadataParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "File.GetMetadata")
}

func TestFileUpload(t *testing.T) {
	client := contractClient(t, "FileService", "Upload")
	response, err := client.Files.Upload(context.Background(), managed.FileUploadParams{File: strings.NewReader("file")})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "File.Upload")
}

func TestPresignedFileAndSkillDownloads(t *testing.T) {
	calls := 0
	c := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		switch calls {
		case 1:
			if r.Header.Get("Authorization") != "Bearer secret-pat" {
				t.Fatal("missing PAT")
			}
			return reply(r, 200, `{"url":"https://storage.test/content?signature=signed","expires_at":"2026-09-09T00:00:00Z"}`), nil
		case 2:
			if r.URL.Host != "storage.test" || r.URL.Query().Get("signature") != "signed" {
				t.Fatal(r.URL)
			}
			if r.Header.Get("Authorization") != "" {
				t.Fatal("PAT leaked to storage")
			}
			return reply(r, 200, "actual file"), nil
		default:
			return reply(r, 200, "ZIP bytes"), nil
		}
	})
	res, e := c.Files.Download(context.Background(), "file", managed.FileDownloadParams{})
	if e != nil {
		t.Fatal(e)
	}
	b, e := io.ReadAll(res.Body)
	res.Body.Close()
	if e != nil || string(b) != "actual file" {
		t.Fatal("download returned link JSON instead of bytes")
	}
	res, e = c.Skills.Versions.Download(context.Background(), "123", managed.SkillVersionDownloadParams{SkillID: "skill"})
	if e != nil {
		t.Fatal(e)
	}
	b, e = io.ReadAll(res.Body)
	res.Body.Close()
	if e != nil || string(b) != "ZIP bytes" || calls != 3 {
		t.Fatal("skill archive was not streamed directly")
	}
}
