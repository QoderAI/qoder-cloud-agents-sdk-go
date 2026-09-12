//go:build live

package managed_test

import (
	"context"
	"errors"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
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
	metadata := liveResult(s.client.Files.GetMetadata(ctx, created.ID, managed.FileGetMetadataParams{})).require(t)
	// Uploads through this SDK carry no purpose, so the API stores them as
	// user_upload and refuses to serve the bytes back.
	if metadata.Downloadable {
		t.Fatalf("expected a non-downloadable upload: %s", metadata.RawJSON())
	}
	_, err := s.client.Files.Download(ctx, created.ID, managed.FileDownloadParams{})
	var apiErr *convention.Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 403 {
		t.Fatalf("expected 403 for downloading an upload, got %v", err)
	}

}
