//go:build live

package forward_test

import (
	"context"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestValidationOnlyBatchLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE", "EXECUTION")
	ctx := s.context(t)
	file := s.file(t, "sdk-validation.jsonl", "session_resource", "{\"custom_id\":\"sdk-live-validation\",\"body\":{\"input\":\"missing required Forward batch fields\"}}\n")
	batch, err := s.client.Batches.New(ctx, forward.BatchNewParams{InputFileID: file.ID, CompletionWindow: "24h"})
	liveCheck(t, err)
	s.cleanup(t, "batch", func(ctx context.Context) error {
		got, err := s.client.Batches.Get(ctx, batch.ID)
		if err != nil {
			return err
		}
		switch got.Status {
		case "completed", "failed", "cancelled", "expired":
			return nil
		}
		_, err = s.client.Batches.Cancel(ctx, batch.ID, forward.BatchCancelParams{})
		return err
	})
	got, err := s.client.Batches.Get(ctx, batch.ID)
	liveCheck(t, err)
	if got.ID != batch.ID {
		t.Fatal("batch did not round trip")
	}
	_, err = s.client.Batches.Tasks.List(ctx, batch.ID, forward.BatchTaskListParams{})
	liveCheck(t, err)
}
