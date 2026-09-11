//go:build live

package forward_test

import (
	"context"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestChannelAndQRSessionLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE", "CHANNEL")
	ctx := s.context(t)
	identity := s.identity(t)
	template := s.template(t, s.environment(t).ID)
	channel, err := s.client.Channels.New(ctx, forward.ChannelNewParams{ChannelType: "wechat", IdentityID: forward.String(identity.ID), TemplateID: forward.String(template.ID), Name: forward.String(liveName("channel"))})
	liveCheck(t, err)
	s.cleanup(t, "channel", func(ctx context.Context) error { _, err := s.client.Channels.Delete(ctx, channel.ID); return err })
	got, err := s.client.Channels.Get(ctx, channel.ID)
	liveCheck(t, err)
	if got.ID != channel.ID {
		t.Fatal("channel did not round trip")
	}
	qr, err := s.client.Channels.QRSessions.New(ctx, channel.ID, forward.ChannelQRSessionNewParams{})
	liveCheck(t, err)
	if qr.SessionKey == "" {
		t.Fatal("missing QR session key")
	}
	gotQR, err := s.client.Channels.QRSessions.Get(ctx, qr.SessionKey)
	liveCheck(t, err)
	if gotQR.SessionKey != qr.SessionKey {
		t.Fatal("QR session did not round trip")
	}
	updated, err := s.client.Channels.Update(ctx, channel.ID, forward.ChannelUpdateParams{Enabled: forward.Bool(false)})
	liveCheck(t, err)
	if updated.Enabled {
		t.Fatal("channel still enabled")
	}
}
