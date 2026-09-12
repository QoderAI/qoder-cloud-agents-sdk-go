//go:build live

package managed_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type managedScenarioSuite struct {
	client                     managed.Client
	timeout                    time.Duration
	scenarioTimeout            time.Duration
	allowWrite, allowExecution bool
}

func newManagedScenarioSuite(t *testing.T) *managedScenarioSuite {
	t.Helper()
	token := strings.TrimSpace(os.Getenv("QODER_MANAGED_PAT"))
	if token == "" {
		t.Skip("QODER_MANAGED_PAT is not configured")
	}
	timeout := 15 * time.Second
	if raw := os.Getenv("QODER_MANAGED_LIVE_TIMEOUT_SECONDS"); raw != "" {
		seconds, err := strconv.Atoi(raw)
		if err != nil || seconds <= 0 {
			t.Fatal("QODER_MANAGED_LIVE_TIMEOUT_SECONDS must be a positive integer")
		}
		timeout = time.Duration(seconds) * time.Second
	}
	base := os.Getenv("QODER_MANAGED_BASE_URL")
	if base == "" {
		base = managed.DefaultBaseURL
	}
	return &managedScenarioSuite{
		client:  managed.NewClient(option.WithAccessToken(token), option.WithBaseURL(base), option.WithRequestTimeout(timeout), option.WithMaxRetries(0), testutil.RequestLog(t)),
		timeout: timeout, scenarioTimeout: testutil.ExecutionTimeout(t), allowWrite: os.Getenv("QODER_MANAGED_LIVE_ALLOW_WRITE") == "true", allowExecution: os.Getenv("QODER_MANAGED_LIVE_ALLOW_EXECUTION") == "true",
	}
}
func (s *managedScenarioSuite) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.scenarioTimeout)
}
func (s *managedScenarioSuite) requireWrite(t *testing.T) {
	t.Helper()
	if !s.allowWrite {
		t.Skip("QODER_MANAGED_LIVE_ALLOW_WRITE is not true")
	}
}
func (s *managedScenarioSuite) requireExecution(t *testing.T) {
	t.Helper()
	s.requireWrite(t)
	if !s.allowExecution {
		t.Skip("QODER_MANAGED_LIVE_ALLOW_EXECUTION is not true")
	}
}
func (s *managedScenarioSuite) cleanup(t *testing.T, label string, run func(context.Context) error) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.scenarioTimeout)
		defer cancel()
		if err := run(ctx); err != nil {
			if testutil.ResourceAlreadyGone(err) {
				return
			}
			t.Errorf("cleanup %s: %s", label, testutil.SafeError(err))
		}
	})
}
func managedUnique(prefix string) string {
	return fmt.Sprintf("sdk-%s-%d", prefix, time.Now().UnixNano())
}

type liveResponse[T any] struct {
	value T
	err   error
}

func liveResult[T any](value T, err error) liveResponse[T] { return liveResponse[T]{value, err} }
func (r liveResponse[T]) require(t *testing.T) T {
	t.Helper()
	if r.err != nil {
		t.Fatal(testutil.SafeError(r.err))
	}
	return r.value
}

func liveJSON[T any](t *testing.T, data any) T {
	t.Helper()
	var value T
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(encoded, &value); err != nil {
		t.Fatal(err)
	}
	return value
}
func (s *managedScenarioSuite) model(t *testing.T) string {
	t.Helper()
	if model := os.Getenv("QODER_MANAGED_MODEL"); model != "" {
		return model
	}
	ctx, cancel := s.context()
	defer cancel()
	models := liveResult(s.client.Models.List(ctx, managed.ModelListParams{})).require(t)
	for _, model := range models.Data {
		if model.IsEnabled && model.ID != "" {
			return model.ID
		}
	}
	t.Skip("the Managed account has no enabled model")
	return ""
}
func (s *managedScenarioSuite) createEnvironment(t *testing.T) *managed.Environment {
	t.Helper()
	return s.createEnvironmentWithConfig(t, map[string]any{"type": "cloud"})
}
func (s *managedScenarioSuite) createEnvironmentWithConfig(t *testing.T, config map[string]any) *managed.Environment {
	t.Helper()
	s.requireWrite(t)
	ctx, cancel := s.context()
	defer cancel()
	params := liveJSON[managed.EnvironmentNewParams](t, map[string]any{"name": managedUnique("env"), "config": config, "metadata": map[string]string{"suite": "sdk-live"}})
	value := liveResult(s.client.Environments.New(ctx, params)).require(t)
	s.cleanup(t, "Environment "+value.ID, func(ctx context.Context) error {
		_, err := s.client.Environments.Delete(ctx, value.ID, managed.EnvironmentDeleteParams{})
		var api *convention.Error
		if errors.As(err, &api) && api.StatusCode == 409 {
			_, err = s.client.Environments.Archive(ctx, value.ID, managed.EnvironmentArchiveParams{})
		}
		return err
	})
	return value
}
func (s *managedScenarioSuite) createAgent(t *testing.T) *managed.ManagedAgentsAgent {
	t.Helper()
	s.requireWrite(t)
	model := s.model(t)
	ctx, cancel := s.context()
	defer cancel()
	value := liveResult(s.client.Agents.New(ctx, managed.AgentNewParams{Name: managedUnique("agent"), Model: managed.ManagedAgentsModelConfigParams{ID: model}, System: managed.String("You are a Managed SDK test assistant."), Metadata: map[string]string{"suite": "sdk-live"}}, option.WithHeader("Idempotency-Key", managedUnique("ikey")))).require(t)
	s.cleanup(t, "Agent "+value.ID, func(ctx context.Context) error {
		_, err := s.client.Agents.Archive(ctx, value.ID, managed.AgentArchiveParams{})
		return err
	})
	return value
}
func (s *managedScenarioSuite) createMemoryStore(t *testing.T) *managed.ManagedAgentsMemoryStore {
	t.Helper()
	s.requireWrite(t)
	ctx, cancel := s.context()
	defer cancel()
	value := liveResult(s.client.MemoryStores.New(ctx, managed.MemoryStoreNewParams{Name: managedUnique("memory"), Metadata: map[string]string{"suite": "sdk-live"}})).require(t)
	s.cleanup(t, "Memory Store "+value.ID, func(ctx context.Context) error {
		_, err := s.client.MemoryStores.Delete(ctx, value.ID, managed.MemoryStoreDeleteParams{})
		return err
	})
	return value
}
