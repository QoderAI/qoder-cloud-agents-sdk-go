//go:build live

package forward_test

import (
	"context"
	"os"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestChannelPairingLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE", "CHANNEL")
	code := os.Getenv("QODER_FORWARD_LIVE_PAIRING_CODE")
	if code == "" {
		t.Fatal("QODER_FORWARD_LIVE_PAIRING_CODE is required")
	}
	identity := s.identity(t)
	template := s.template(t, s.environment(t).ID)
	pairing, err := s.client.ChannelPairings.New(s.context(t), forward.ChannelPairingNewParams{
		Code:           code,
		IdentityID:     identity.ID,
		TemplateID:     template.ID,
		IdempotencyKey: forward.String(liveName("pairing")),
	})
	liveCheck(t, err)
	if pairing.ID == "" || pairing.Status != "active" {
		t.Fatalf("unexpected pairing response: id=%q status=%q", pairing.ID, pairing.Status)
	}
	if want := os.Getenv("QODER_FORWARD_LIVE_CHANNEL_ID"); want != "" && pairing.ChannelID != want {
		t.Fatalf("channel ID = %q, want %q", pairing.ChannelID, want)
	}
	deleted := false
	s.cleanup(t, "channel pairing", func(ctx context.Context) error {
		if deleted {
			return nil
		}
		_, err := s.client.ChannelPairings.Delete(ctx, pairing.ID)
		return err
	})
	result, err := s.client.ChannelPairings.Delete(s.context(t), pairing.ID)
	liveCheck(t, err)
	if result.ID != pairing.ID || !result.Deleted {
		t.Fatalf("unexpected unpair response: id=%q deleted=%t", result.ID, result.Deleted)
	}
	deleted = true
}
