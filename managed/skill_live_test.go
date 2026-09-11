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

func TestSkillListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.Skills.List(ctx, managed.SkillListParams{})).require(t)
}

func TestSkillLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()

	s.requireWrite(t)
	name := managedUnique("skill")
	markdown := "---\nname: " + name + "\ndescription: Managed SDK live skill\n---\n# " + name + "\nReply with SDK-LIVE.\n"
	files := func() []io.Reader {
		return []io.Reader{convention.UploadFile{Reader: strings.NewReader(markdown), Name: name + "/SKILL.md", MediaType: "text/markdown"}}
	}
	created := liveResult(s.client.Skills.New(ctx, managed.SkillNewParams{Files: files(), DisplayName: managed.String(name), Metadata: map[string]string{"suite": "sdk-live"}})).require(t)
	s.cleanup(t, "Skill "+created.ID, func(ctx context.Context) error {
		_, err := s.client.Skills.Delete(ctx, created.ID, managed.SkillDeleteParams{})
		return err
	})
	liveResult(s.client.Skills.Get(ctx, created.ID, managed.SkillGetParams{})).require(t)
	version := liveResult(s.client.Skills.Versions.New(ctx, created.ID, managed.SkillVersionNewParams{Files: files()})).require(t)
	s.cleanup(t, "Skill version "+version.Version, func(ctx context.Context) error {
		_, err := s.client.Skills.Versions.Delete(ctx, version.Version, managed.SkillVersionDeleteParams{SkillID: created.ID})
		return err
	})
	liveResult(s.client.Skills.Versions.Get(ctx, version.Version, managed.SkillVersionGetParams{SkillID: created.ID})).require(t)
	response := liveResult(s.client.Skills.Versions.Download(ctx, version.Version, managed.SkillVersionDownloadParams{SkillID: created.ID})).require(t)
	defer response.Body.Close()
	if len(liveResult(io.ReadAll(response.Body)).require(t)) == 0 {
		t.Fatal("empty skill archive")
	}

}
