//go:build live

package forward_test

import (
	"io"
	"testing"
)

func TestFileUploadDownloadLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	// Only tool_output, skill_output, session_resource and agent_output are
	// downloadable; a user_upload file is write-only by design.
	file := s.file(t, "sdk-live.txt", "session_resource", "SDK live file content")
	got, err := s.client.Files.GetMetadata(ctx, file.ID)
	liveCheck(t, err)
	if got.ID != file.ID {
		t.Fatal("file metadata did not round trip")
	}
	if !got.Downloadable {
		t.Fatalf("file is not downloadable: %s", got.RawJSON())
	}
	response, err := s.client.Files.Download(ctx, file.ID)
	liveCheck(t, err)
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	liveCheck(t, err)
	if string(data) != "SDK live file content" {
		t.Fatal("downloaded content differs")
	}
}
