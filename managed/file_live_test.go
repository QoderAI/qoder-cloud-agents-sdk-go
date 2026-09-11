//go:build live

package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"io"
	"strings"
	"testing"
)

func TestFileListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.Files.List(ctx, managed.FileListParams{})).require(t)
}

func TestFileLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()

	s.requireWrite(t)
	content := "Managed SDK live file\n"
	created := liveResult(s.client.Files.Upload(ctx, managed.FileUploadParams{File: convention.UploadFile{Reader: strings.NewReader(content), Name: managedUnique("file") + ".txt"}, Metadata: map[string]string{"suite": "sdk-live"}})).require(t)
	s.cleanup(t, "File "+created.ID, func(ctx context.Context) error {
		_, err := s.client.Files.Delete(ctx, created.ID, managed.FileDeleteParams{})
		return err
	})
	liveResult(s.client.Files.GetMetadata(ctx, created.ID, managed.FileGetMetadataParams{})).require(t)
	response := liveResult(s.client.Files.Download(ctx, created.ID, managed.FileDownloadParams{})).require(t)
	defer response.Body.Close()
	if got := string(liveResult(io.ReadAll(response.Body)).require(t)); got != content {
		t.Fatalf("download content: %q", got)
	}

}
