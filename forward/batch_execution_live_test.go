//go:build live

package forward_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestForwardBatchE2ELive(t *testing.T) {
	testutil.RequireE2E(t, "FORWARD")
	s := newLiveSuite(t, "WRITE", "EXECUTION")
	ctx, cancel := context.WithTimeout(context.Background(), testutil.ExecutionTimeout(t))
	defer cancel()
	identity := s.identity(t)
	template := s.template(t, s.environment(t).ID)
	marker := testutil.Marker(t)
	customID := liveName("task")
	line, err := json.Marshal(map[string]any{"custom_id": customID, "template_id": template.ID, "identity_id": identity.ID, "body": map[string]any{"input": "Reply with exactly " + marker}})
	liveCheck(t, err)
	input := s.file(t, "sdk-e2e-input.jsonl", "session_resource", string(line)+"\n")
	// The scheduler only promotes queued batches inside the server's idle
	// window (22:00–08:00) unless the batch opts out.
	batch, err := s.client.Batches.New(ctx, forward.BatchNewParams{InputFileID: input.ID, CompletionWindow: "24h", IdempotencyKey: forward.String(liveName("batch"))}, option.WithJSONSet("ignore_idle_window", true))
	liveCheck(t, err)
	t.Logf("batch=%s input_file=%s custom_id=%s", batch.ID, input.ID, customID)
	batchID := batch.ID
	s.cleanup(t, "batch "+batchID, func(ctx context.Context) error {
		return s.finishBatch(ctx, batchID, customID, identity.ID, template.ID)
	})
	for !batchTerminal(batch.Status) {
		batch, err = s.client.Batches.Get(ctx, batchID)
		liveCheck(t, err)
		if batchTerminal(batch.Status) {
			break
		}
		if err = testutil.PollPause(ctx); err != nil {
			t.Fatalf("batch=%s status=%s (execution depends on the server's batch window): %v", batchID, batch.Status, err)
		}
	}
	if batch.Status != "completed" || batch.RequestCounts.Completed != 1 || batch.RequestCounts.Failed != 0 || batch.OutputFileID == "" {
		t.Fatalf("batch=%s status=%s counts=%+v output=%s", batchID, batch.Status, batch.RequestCounts, batch.OutputFileID)
	}
	tasks, err := s.client.Batches.Tasks.List(ctx, batchID, forward.BatchTaskListParams{})
	liveCheck(t, err)
	if len(tasks.Data) != 1 || tasks.Data[0].CustomID != customID {
		t.Fatal("batch task did not round trip")
	}
	rows, err := s.batchOutput(ctx, batchID)
	liveCheck(t, err)
	if len(rows) != 1 {
		t.Fatalf("output rows=%d want=1", len(rows))
	}
	row := rows[0]
	if row.CustomID != customID || row.IdentityID != identity.ID || row.TemplateID != template.ID || row.Status != "completed" || row.SessionID == "" || string(row.Error) != "null" && len(row.Error) != 0 {
		t.Fatal("unexpected batch result row")
	}
	s.waitTurn(t, row.SessionID, "", []string{marker}, false, false)
	if !strings.Contains(string(row.Response), marker) {
		t.Fatal("batch output did not contain the assistant result")
	}

}
func batchTerminal(status string) bool {
	return status == "completed" || status == "failed" || status == "cancelled" || status == "expired"
}

type batchOutputRow struct {
	CustomID   string          `json:"custom_id"`
	SessionID  string          `json:"session_id"`
	TemplateID string          `json:"template_id"`
	IdentityID string          `json:"identity_id"`
	Status     string          `json:"status"`
	Response   json.RawMessage `json:"response"`
	Error      json.RawMessage `json:"error"`
}

func (s *liveSuite) batchOutput(ctx context.Context, id string) ([]batchOutputRow, error) {
	link, err := s.client.Batches.GetOutput(ctx, id)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, link.URL, nil)
	if err != nil {
		return nil, err
	}
	output, err := (&http.Client{Timeout: s.timeout}).Do(request)
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()
	if output.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("batch output HTTP %d", output.StatusCode)
	}
	scanner := bufio.NewScanner(io.LimitReader(output.Body, 4<<20))
	scanner.Buffer(make([]byte, 4096), 2<<20)
	var rows []batchOutputRow
	for scanner.Scan() {
		var row batchOutputRow
		if err = json.Unmarshal(scanner.Bytes(), &row); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, scanner.Err()
}
func (s *liveSuite) finishBatch(ctx context.Context, id, customID, identityID, templateID string) error {
	current, err := s.client.Batches.Get(ctx, id)
	if err != nil {
		return err
	}
	if !batchTerminal(current.Status) {
		if _, err = s.client.Batches.Cancel(ctx, id, forward.BatchCancelParams{}); err != nil {
			return err
		}
	}
	for !batchTerminal(current.Status) {
		current, err = s.client.Batches.Get(ctx, id)
		if err != nil {
			return err
		}
		if batchTerminal(current.Status) {
			break
		}
		if err = testutil.PollPause(ctx); err != nil {
			return fmt.Errorf("batch=%s remains %s: %w", id, current.Status, err)
		}
	}
	if current.OutputFileID == "" {
		if current.RequestCounts.Total == 0 {
			return nil
		}
		return fmt.Errorf("batch=%s has no output for session cleanup", id)
	}
	rows, err := s.batchOutput(ctx, id)
	if err != nil {
		return &testutil.CleanupFailure{Err: err}
	}
	if len(rows) != 1 || rows[0].CustomID != customID || rows[0].IdentityID != identityID || rows[0].TemplateID != templateID {
		return fmt.Errorf("batch=%s cleanup output does not match the test input", id)
	}
	if rows[0].SessionID != "" {
		return s.finishSession(ctx, rows[0].SessionID)
	}
	return nil
}
