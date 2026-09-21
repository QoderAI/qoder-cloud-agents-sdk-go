//go:build live

package forward_test

import (
	"context"
	"encoding/json"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestScheduleCleanupOffline(t *testing.T) {
	var requests []string
	runReads := 0
	sessionReads := 0
	s := &liveSuite{client: testClient(func(r *http.Request) (*http.Response, error) {
		key := r.Method + " " + r.URL.Path
		requests = append(requests, key)
		switch key {
		case "GET /api/v1/forward/schedule_runs/run":
			runReads++
			if runReads == 1 {
				return reply(r, 200, `{"id":"run","status":"pending"}`), nil
			}
			return reply(r, 200, `{"id":"run","status":"running","session_id":"session"}`), nil
		case "GET /api/v1/forward/sessions/session":
			sessionReads++
			status := "running"
			if sessionReads > 1 {
				status = "idle"
			}
			return reply(r, 200, `{"id":"session","status":"`+status+`"}`), nil
		case "POST /api/v1/forward/sessions/session/cancel":
			return reply(r, 200, `{"id":"session"}`), nil
		case "POST /api/v1/forward/sessions/session/archive":
			return reply(r, 200, `{"id":"session"}`), nil
		default:
			t.Fatalf("unexpected request %s", key)
			return nil, nil
		}
	})}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.finishScheduleRun(ctx, "run", "identity"); err != nil {
		t.Fatal(err)
	}
	want := []string{"GET /api/v1/forward/schedule_runs/run", "GET /api/v1/forward/schedule_runs/run", "GET /api/v1/forward/sessions/session", "POST /api/v1/forward/sessions/session/cancel", "GET /api/v1/forward/sessions/session", "POST /api/v1/forward/sessions/session/archive"}
	if !reflect.DeepEqual(requests, want) {
		t.Fatalf("cleanup order=%v", requests)
	}
}

func TestPendingScheduleCleanupDeadlineOffline(t *testing.T) {
	s := &liveSuite{client: testClient(func(r *http.Request) (*http.Response, error) {
		return reply(r, 200, `{"id":"run","status":"pending"}`), nil
	})}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if err := s.finishScheduleRun(ctx, "run", "identity"); err == nil || !strings.Contains(err.Error(), "run=run") {
		t.Fatalf("pending run disappeared from diagnostics: %v", err)
	}
}

func TestBatchFailureCleanupOffline(t *testing.T) {
	for _, state := range []string{"failed", "processing"} {
		t.Run(state, func(t *testing.T) {
			storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "" {
					t.Error("API token leaked to batch storage")
				}
				_, _ = w.Write([]byte(`{"custom_id":"task","identity_id":"identity","template_id":"template","session_id":"session","status":"failed","error":{"code":"timeout"}}` + "\n"))
			}))
			defer storage.Close()
			var requests []string
			batchReads := 0
			s := &liveSuite{timeout: time.Second, client: testClient(func(r *http.Request) (*http.Response, error) {
				key := r.Method + " " + r.URL.Path
				requests = append(requests, key)
				switch key {
				case "GET /api/v1/forward/batches/batch":
					batchReads++
					status := state
					if batchReads > 1 {
						status = "cancelled"
					}
					return reply(r, 200, `{"id":"batch","status":"`+status+`","output_file_id":"internal-file","request_counts":{"total":1}}`), nil
				case "POST /api/v1/forward/batches/batch/cancel":
					return reply(r, 200, `{"id":"batch","status":"cancelling"}`), nil
				case "GET /api/v1/forward/batches/batch/output":
					return reply(r, 200, `{"url":"`+storage.URL+`/output"}`), nil
				case "GET /api/v1/forward/sessions/session":
					return reply(r, 200, `{"id":"session","status":"idle"}`), nil
				case "POST /api/v1/forward/sessions/session/archive":
					return reply(r, 200, `{"id":"session"}`), nil
				default:
					t.Fatalf("unexpected request %s", key)
					return nil, nil
				}
			})}
			if err := s.finishBatch(context.Background(), "batch", "task", "identity", "template"); err != nil {
				t.Fatal(err)
			}
			if got := requests[len(requests)-1]; got != "POST /api/v1/forward/sessions/session/archive" {
				t.Fatalf("failed batch session not archived: %v", requests)
			}
			if strings.Contains(strings.Join(requests, "\n"), "/files/") {
				t.Fatal("internal output treated as a public file")
			}
		})
	}
}

func TestBatchMissingOutputCleanupOffline(t *testing.T) {
	s := &liveSuite{client: testClient(func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "/output") {
			return reply(r, 404, `{"error":{"message":"output not generated"}}`), nil
		}
		return reply(r, 200, `{"id":"batch","status":"cancelled","output_file_id":"internal","request_counts":{"total":1}}`), nil
	})}
	err := s.finishBatch(context.Background(), "batch", "task", "identity", "template")
	if err == nil || testsupport.ResourceAlreadyGone(err) {
		t.Fatalf("missing output silently ignored: %v", err)
	}
}

// A child test deliberately fails its first poll. The parent verifies that Go
// still executes the registered cleanup and that it uses the immutable run ID.
func TestExecutionFailureCleanupOffline(t *testing.T) {
	for _, scenario := range []string{"schedule", "batch"} {
		t.Run(scenario, func(t *testing.T) {
			var mu sync.Mutex
			var requests []string
			polls := 0
			customID := ""
			var baseURL string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				key := r.Method + " " + r.URL.Path
				requests = append(requests, key)
				w.Header().Set("Content-Type", "application/json")
				respond := func(body string) { _, _ = io.WriteString(w, body) }
				switch key {
				case "POST /api/v1/forward/environments":
					respond(`{"id":"env"}`)
				case "POST /api/v1/forward/identities":
					respond(`{"id":"identity"}`)
				case "POST /api/v1/forward/templates":
					respond(`{"id":"template"}`)
				case "POST /api/v1/forward/schedules":
					respond(`{"id":"schedule"}`)
				case "POST /api/v1/forward/schedules/schedule/run":
					respond(`{"id":"run","status":"pending"}`)
				case "POST /api/v1/forward/files":
					if err := r.ParseMultipartForm(1 << 20); err != nil {
						t.Error(err)
						w.WriteHeader(400)
						return
					}
					defer r.MultipartForm.RemoveAll()
					file, _, err := r.FormFile("file")
					if err != nil {
						t.Error(err)
						w.WriteHeader(400)
						return
					}
					defer file.Close()
					var row struct {
						CustomID string `json:"custom_id"`
					}
					if err = json.NewDecoder(file).Decode(&row); err != nil {
						t.Error(err)
					}
					customID = row.CustomID
					respond(`{"id":"input"}`)
				case "POST /api/v1/forward/batches":
					respond(`{"id":"batch","status":"processing"}`)
				case "GET /api/v1/forward/schedule_runs/run", "GET /api/v1/forward/batches/batch":
					polls++
					if polls == 1 {
						w.WriteHeader(503)
						respond(`{"request_id":"failed-poll","error":{"type":"api_error","message":"deliberate test failure"}}`)
						return
					}
					if scenario == "schedule" {
						respond(`{"id":"run","status":"running","session_id":"session"}`)
					} else {
						respond(`{"id":"batch","status":"cancelled","output_file_id":"internal","request_counts":{"total":1}}`)
					}
				case "GET /api/v1/forward/batches/batch/output":
					respond(`{"url":"` + baseURL + `/output"}`)
				case "GET /output":
					respond(`{"custom_id":"` + customID + `","identity_id":"identity","template_id":"template","session_id":"session","status":"failed"}` + "\n")
				case "GET /api/v1/forward/sessions/session":
					respond(`{"id":"session","status":"idle"}`)
				case "POST /api/v1/forward/sessions/session/archive", "POST /api/v1/forward/schedules/schedule/archive", "POST /api/v1/forward/templates/template/archive":
					respond(`{}`)
				case "DELETE /api/v1/forward/environments/env", "DELETE /api/v1/forward/identities/identity", "DELETE /api/v1/forward/files/input":
					respond(`{}`)
				default:
					t.Errorf("unexpected local request %s", key)
					w.WriteHeader(400)
				}
			}))
			defer server.Close()
			mu.Lock()
			baseURL = server.URL
			mu.Unlock()
			command := exec.Command(os.Args[0], "-test.run=^TestExecutionFailureCleanupChild$", "-test.timeout=15s")
			command.Env = append(os.Environ(), "QODER_SDK_CLEANUP_CHILD="+scenario, "QODER_SDK_CLEANUP_BASE="+baseURL)
			output, err := command.CombinedOutput()
			if err == nil || !strings.Contains(string(output), "failed-poll") || strings.Contains(string(output), "panic:") {
				t.Fatalf("unexpected child result: %v\n%s", err, output)
			}
			mu.Lock()
			defer mu.Unlock()
			if !slices.Contains(requests, "POST /api/v1/forward/sessions/session/archive") {
				t.Fatalf("session cleanup missing after failed poll: %v\n%s", requests, output)
			}
			if polls != 2 {
				t.Fatalf("cleanup did not resolve stable ID: polls=%d", polls)
			}
		})
	}
}
func TestExecutionFailureCleanupChild(t *testing.T) {
	scenario := os.Getenv("QODER_SDK_CLEANUP_CHILD")
	if scenario == "" {
		t.Skip("subprocess fixture")
	}
	t.Setenv("QODER_FORWARD_BASE_URL", os.Getenv("QODER_SDK_CLEANUP_BASE")+"/api/v1/forward")
	t.Setenv("QODER_FORWARD_PAT", "offline-test-token")
	t.Setenv("QODER_FORWARD_MODEL", "offline-model")
	t.Setenv("QODER_FORWARD_LIVE_ALLOW_WRITE", "true")
	t.Setenv("QODER_FORWARD_LIVE_ALLOW_EXECUTION", "true")
	t.Setenv("QODER_E2E_TIMEOUT_SECONDS", "5")
	if scenario == "schedule" {
		TestForwardScheduleE2ELive(t)
	} else {
		TestForwardBatchE2ELive(t)
	}
}
