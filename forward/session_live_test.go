//go:build live

package forward_test

import (
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestSessionResourceAndThreadLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	identity := s.identity(t)
	template := s.template(t, s.environment(t).ID)
	session := s.session(t, identity.ID, template.ID)
	updated, err := s.client.Sessions.Update(ctx, session.ID, forward.SessionUpdateParams{Title: forward.String("SDK updated session")})
	liveCheck(t, err)
	if updated.Title != "SDK updated session" {
		t.Fatal("session title not updated")
	}
	file := s.file(t, "sdk-resource.txt", "session_resource", "SDK resource")
	resource, err := s.client.Sessions.Resources.Add(ctx, session.ID, forward.SessionResourceAddParams{Type: "file", FileID: file.ID, MountPath: forward.String("/data/workspace/sdk-resource.txt")})
	liveCheck(t, err)
	if resource.FileID != file.ID || resource.MountPath != "/data/workspace/sdk-resource.txt" {
		t.Fatal("resource mount did not round trip")
	}
	_, err = s.client.Sessions.Events.List(ctx, session.ID, forward.SessionEventListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	threads, err := s.client.Sessions.Threads.List(ctx, session.ID, forward.SessionThreadListParams{})
	liveCheck(t, err)
	for _, thread := range threads.Data {
		got, err := s.client.Sessions.Threads.Get(ctx, session.ID, thread.ID)
		liveCheck(t, err)
		if got.ID != thread.ID {
			t.Fatal("thread did not round trip")
		}
		_, err = s.client.Sessions.Threads.Events.List(ctx, session.ID, thread.ID, forward.SessionThreadEventListParams{})
		liveCheck(t, err)
	}
}
