package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/forwardutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	config, err := live.Load("forward")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	s := &forwardutil.Example{Client: forward.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "batch",
		Description: "提交一个 JSONL 批处理任务，等待完成后检查结果。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

// run bypasses the idle window so the example can run at any time of day.
func run(ctx context.Context, r *live.Run, s *forwardutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	identity, err := s.Identity(ctx, r)
	if err != nil {
		return err
	}
	template, err := s.Template(ctx, r, forward.TemplateNewParams{EnvironmentID: env})
	if err != nil {
		return err
	}
	marker, customID := live.Marker(), live.Name("task")
	r.Log("info", "批处理任务的输入消息：")
	r.Message("user", "Reply with exactly "+marker)
	line, err := json.Marshal(map[string]any{"custom_id": customID, "template_id": template, "identity_id": identity, "body": map[string]any{"input": "Reply with exactly " + marker}})
	if err != nil {
		return err
	}
	r.Step("上传包含一个任务的 JSONL 输入文件")
	input, err := s.Client.Files.Upload(ctx, forward.FileUploadParams{File: convention.UploadFile{Name: "sdk-example-input.jsonl", Reader: strings.NewReader(string(line) + "\n")}, Purpose: forward.String("session_resource")})
	if err != nil {
		return err
	}
	r.Track("input_file", input.ID, func(ctx context.Context) error { return s.Client.Files.Delete(ctx, input.ID) })
	r.Step("提交 Batch，跳过闲时窗口限制")
	batch, err := s.Client.Batches.New(ctx, forward.BatchNewParams{
		InputFileID:      input.ID,
		CompletionWindow: "24h",
		IdempotencyKey:   forward.String(live.Name("batch")),
	}, option.WithJSONSet("ignore_idle_window", true))
	if err != nil {
		return err
	}
	batchID := batch.ID
	r.Track("batch", batchID, func(ctx context.Context) error { return finishBatch(ctx, s, r, batchID, customID, identity, template) })
	stopWaiting := r.Wait("等待 Batch 完成")
	defer stopWaiting()
	previous := ""
	for !batchTerminal(batch.Status) {
		if previous != batch.Status {
			r.Log("info", "Batch 状态："+live.Status(batch.Status))
			previous = batch.Status
		}
		if err = live.Pause(ctx); err != nil {
			return fmt.Errorf("batch=%s status=%s: %w", batchID, previous, err)
		}
		batch, err = s.Client.Batches.Get(ctx, batchID)
		if err != nil {
			return err
		}
	}
	stopWaiting()
	if batch.Status != "completed" || batch.RequestCounts.Completed != 1 || batch.RequestCounts.Failed != 0 || batch.OutputFileID == "" {
		return fmt.Errorf("batch=%s status=%s completed=%d failed=%d", batchID, batch.Status, batch.RequestCounts.Completed, batch.RequestCounts.Failed)
	}
	r.Step("检查 Batch 的任务和输出文件")
	tasks, err := s.Client.Batches.Tasks.List(ctx, batchID, forward.BatchTaskListParams{})
	if err != nil {
		return err
	}
	if len(tasks.Data) != 1 || tasks.Data[0].CustomID != customID {
		return errors.New("batch task did not round trip")
	}
	rows, err := batchOutput(ctx, s, batchID)
	if err != nil {
		return err
	}
	if len(rows) != 1 {
		return errors.New("expected one batch output row")
	}
	row := rows[0]
	if row.CustomID != customID || row.IdentityID != identity || row.TemplateID != template || row.SessionID == "" || row.Status != "completed" || (len(row.Error) != 0 && string(row.Error) != "null") {
		return errors.New("batch output ownership/status mismatch")
	}
	if !strings.Contains(string(row.Response), marker) {
		return errors.New("batch result missing expected output")
	}
	return s.WaitTurn(ctx, r, row.SessionID, "", []string{marker}, false, false)
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

func batchOutput(ctx context.Context, s *forwardutil.Example, id string) ([]batchOutputRow, error) {
	link, err := s.Client.Batches.GetOutput(ctx, id)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, link.URL, nil)
	if err != nil {
		return nil, err
	}
	output, err := (&http.Client{Timeout: 30 * time.Second}).Do(request)
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
func finishBatch(ctx context.Context, s *forwardutil.Example, r *live.Run, id, customID, identityID, templateID string) error {
	current, err := s.Client.Batches.Get(ctx, id)
	if err != nil {
		return err
	}
	if !batchTerminal(current.Status) {
		if _, err = s.Client.Batches.Cancel(ctx, id, forward.BatchCancelParams{}); err != nil {
			return err
		}
	}
	for !batchTerminal(current.Status) {
		current, err = s.Client.Batches.Get(ctx, id)
		if err != nil {
			return err
		}
		if batchTerminal(current.Status) {
			break
		}
		if err = live.Pause(ctx); err != nil {
			return fmt.Errorf("batch=%s remains %s: %w", id, current.Status, err)
		}
	}
	r.Log("info", "Batch 最终状态："+live.Status(current.Status))
	if current.OutputFileID == "" {
		if current.RequestCounts.Total == 0 {
			return nil
		}
		return fmt.Errorf("batch=%s has no output for session cleanup", id)
	}
	rows, err := batchOutput(ctx, s, id)
	if err != nil {
		return fmt.Errorf("cannot discover batch sessions: %w", err)
	}
	if len(rows) != 1 || rows[0].CustomID != customID || rows[0].IdentityID != identityID || rows[0].TemplateID != templateID {
		return fmt.Errorf("batch=%s cleanup output does not match the test input", id)
	}
	if rows[0].SessionID != "" {
		return s.FinishSession(ctx, rows[0].SessionID)
	}
	return nil
}
