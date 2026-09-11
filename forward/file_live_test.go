//go:build live

package forward_test

import (
	"io"
	"testing"
)

func TestFileUploadDownloadLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	file := s.file(t, "sdk-live.txt", "user_upload", "SDK live file content")
	got, err := s.client.Files.GetMetadata(ctx, file.ID)
	liveCheck(t, err)
	if got.ID != file.ID {
		t.Fatal("file metadata did not round trip")
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
