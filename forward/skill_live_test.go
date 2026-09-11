//go:build live

package forward_test

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestSkillAndVersionLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	name := liveName("skill")
	files := func(content string) []io.Reader {
		return []io.Reader{convention.UploadFile{Name: name + "/SKILL.md", Reader: strings.NewReader(fmt.Sprintf("---\nname: %s\ndescription: SDK live skill\n---\n# SDK\n%s\n", name, content))}}
	}
	skill, err := s.client.Skills.New(ctx, forward.SkillNewParams{Files: files("version one"), Metadata: map[string]any{"suite": "sdk-live"}})
	liveCheck(t, err)
	s.cleanup(t, "skill", func(ctx context.Context) error { return s.client.Skills.Delete(ctx, skill.ID) })
	if skill.ID == "" || skill.LatestVersion == "" {
		t.Fatal("missing typed Skill response fields")
	}
	got, err := s.client.Skills.Get(ctx, skill.ID, forward.SkillGetParams{})
	liveCheck(t, err)
	if got.ID != skill.ID {
		t.Fatal("skill did not round trip")
	}
	version, err := s.client.Skills.Versions.New(ctx, skill.ID, forward.SkillVersionNewParams{Files: files("version two")})
	liveCheck(t, err)
	if version.Version == "" {
		t.Fatal("missing version number")
	}
	gotVersion, err := s.client.Skills.Versions.Get(ctx, skill.ID, version.Version)
	liveCheck(t, err)
	if gotVersion.ID != version.ID {
		t.Fatal("version did not round trip")
	}
	response, err := s.client.Skills.Versions.Download(ctx, skill.ID, version.Version)
	liveCheck(t, err)
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	liveCheck(t, err)
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	liveCheck(t, err)
	if len(archive.File) == 0 {
		t.Fatal("empty skill archive")
	}
}
