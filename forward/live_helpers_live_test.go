//go:build live

package forward_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

type liveSuite struct {
	client          forward.Client
	timeout         time.Duration
	scenarioTimeout time.Duration
}

func newLiveSuite(t *testing.T, gates ...string) *liveSuite {
	t.Helper()
	if os.Getenv("QODER_FORWARD_PAT") == "" {
		t.Skip("QODER_FORWARD_PAT is not configured")
	}
	for _, gate := range gates {
		if !liveEnabled(gate) {
			t.Skip("QODER_FORWARD_LIVE_ALLOW_" + gate + " is not true")
		}
	}
	seconds := 15
	if s := os.Getenv("QODER_FORWARD_LIVE_TIMEOUT_SECONDS"); s != "" {
		var err error
		seconds, err = strconv.Atoi(s)
		if err != nil || seconds <= 0 {
			t.Fatal("QODER_FORWARD_LIVE_TIMEOUT_SECONDS must be positive")
		}
	}
	timeout := time.Duration(seconds) * time.Second
	return &liveSuite{timeout: timeout, scenarioTimeout: testutil.ExecutionTimeout(t), client: forward.NewClient(option.WithAccessToken(os.Getenv("QODER_FORWARD_PAT")), option.WithRequestTimeout(timeout), option.WithMaxRetries(0), testutil.RequestLog(t))}
}

func liveEnabled(gate string) bool {
	b, _ := strconv.ParseBool(os.Getenv("QODER_FORWARD_LIVE_ALLOW_" + gate))
	return b
}
func liveName(prefix string) string { return fmt.Sprintf("sdk-%s-%d", prefix, time.Now().UnixNano()) }
func liveCheck(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(testutil.SafeError(err))
	}
}
func (s *liveSuite) context(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), s.scenarioTimeout)
	t.Cleanup(cancel)
	return ctx
}
func (s *liveSuite) cleanup(t *testing.T, label string, fn func(context.Context) error) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.scenarioTimeout)
		defer cancel()
		if err := fn(ctx); err != nil {
			if testutil.ResourceAlreadyGone(err) {
				return
			}
			t.Errorf("cleanup %s: %s", label, testutil.SafeError(err))
		}
	})
}
func (s *liveSuite) environment(t *testing.T) *forward.Environment {
	t.Helper()
	r, err := s.client.Environments.New(s.context(t), forward.EnvironmentNewParams{Name: liveName("environment"), Config: map[string]any{"type": "cloud"}})
	liveCheck(t, err)
	s.cleanup(t, "environment", func(ctx context.Context) error {
		err := s.client.Environments.Delete(ctx, r.ID)
		var api *convention.Error
		if errors.As(err, &api) && api.StatusCode == 409 {
			_, err = s.client.Environments.Archive(ctx, r.ID)
		}
		return err
	})
	return r
}
func (s *liveSuite) identity(t *testing.T) *forward.Identity {
	t.Helper()
	r, err := s.client.Identities.New(s.context(t), forward.IdentityNewParams{ExternalID: liveName("identity"), Metadata: map[string]any{"suite": "sdk-live"}})
	liveCheck(t, err)
	s.cleanup(t, "identity", func(ctx context.Context) error { _, err := s.client.Identities.Delete(ctx, r.ID); return err })
	return r
}
func (s *liveSuite) template(t *testing.T, environmentID string) *forward.Template {
	t.Helper()
	model := os.Getenv("QODER_FORWARD_MODEL")
	if model == "" {
		models, err := s.client.Models.List(s.context(t))
		liveCheck(t, err)
		for _, item := range models.Data {
			if item.IsEnabled && item.ID != "" {
				model = item.ID
				break
			}
		}
		if model == "" {
			t.Skip("the live account has no enabled model")
		}
	}
	r, err := s.client.Templates.New(s.context(t), forward.TemplateNewParams{Name: liveName("template"), Model: forward.ModelConfigUnionParam{OfString: forward.String(model)}, EnvironmentID: environmentID, System: forward.String("You are a test assistant.")})
	liveCheck(t, err)
	s.cleanup(t, "template", func(ctx context.Context) error {
		_, err := s.client.Templates.Archive(ctx, r.ID, forward.TemplateArchiveParams{})
		return err
	})
	return r
}
func (s *liveSuite) session(t *testing.T, identityID, templateID string) *forward.Session {
	t.Helper()
	r, err := s.client.Sessions.New(s.context(t), forward.SessionNewParams{IdentityID: identityID, TemplateID: templateID, Title: forward.String("SDK live session")})
	liveCheck(t, err)
	s.cleanupSession(t, r.ID)
	return r
}
func (s *liveSuite) file(t *testing.T, name, purpose, contents string) *forward.FileMetadata {
	t.Helper()
	r, err := s.client.Files.Upload(s.context(t), forward.FileUploadParams{File: convention.UploadFile{Reader: strings.NewReader(contents), Name: name}, Purpose: forward.String(purpose)})
	liveCheck(t, err)
	s.cleanup(t, "file", func(ctx context.Context) error { return s.client.Files.Delete(ctx, r.ID) })
	return r
}
