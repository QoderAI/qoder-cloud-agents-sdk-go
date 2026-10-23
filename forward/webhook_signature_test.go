package forward

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestWebhookSignatureVector(t *testing.T) {
	raw, err := os.ReadFile("testdata/webhook-signature.json")
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]string
	if err = json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	headers := http.Header{"webhook-id": {v["id"]}, "Webhook-Timestamp": {v["timestamp"]}, "Webhook-Signature": {"v2,unknown v1,invalid " + v["signature"]}}
	for _, delta := range []int64{-300, 0, 300} {
		if err = verifyWebhookSignatureAt([]byte(v["payload"]), headers, v["secret"], 1767225600+delta); err != nil {
			t.Fatal(err)
		}
	}
	for _, delta := range []int64{-301, 301} {
		if verifyWebhookSignatureAt([]byte(v["payload"]), headers, v["secret"], 1767225600+delta) == nil {
			t.Fatal("stale accepted")
		}
	}
	for _, mutate := range []func(http.Header){func(h http.Header) { h.Set("Webhook-Timestamp", "9223372036854775807") }, func(h http.Header) { h.Set("Webhook-Timestamp", "-1") }, func(h http.Header) { h.Set("Webhook-Signature", "v2,"+strings.TrimPrefix(v["signature"], "v1,")) }, func(h http.Header) { h.Del("Webhook-Signature") }, func(h http.Header) { h["webhook-id"] = []string{"other"} }} {
		h := headers.Clone()
		mutate(h)
		if verifyWebhookSignatureAt([]byte(v["payload"]), h, v["secret"], 1767225600) == nil {
			t.Fatal("invalid accepted")
		}
	}
	if verifyWebhookSignatureAt([]byte(v["payload"]+" "), headers, v["secret"], 1767225600) == nil {
		t.Fatal("tamper accepted")
	}
	for _, secret := range []string{"", "whsec_", "whsec_notbase64", "wrong_" + v["secret"]} {
		if verifyWebhookSignatureAt([]byte(v["payload"]), headers, secret, 1767225600) == nil {
			t.Fatal("secret accepted")
		}
	}
}
