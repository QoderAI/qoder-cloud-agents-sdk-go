package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestReceiverPreservesUnknownAndSyntheticEvents(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	secret := "whsec_" + base64.StdEncoding.EncodeToString(key)
	for _, body := range []string{`{"type":"event","data":{"type":"future.event","extra":42}}`, `{"type":"event","endpoint_id":"id","note":"synthetic"}`} {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		mac := hmac.New(sha256.New, key)
		mac.Write([]byte("whe_test." + timestamp + "." + body))
		headers := http.Header{"Webhook-Id": {"whe_test"}, "Webhook-Timestamp": {timestamp}, "Webhook-Signature": {"v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))}}
		event, err := parseDelivery([]byte(body), headers, secret)
		if err != nil || event["type"] != "event" {
			t.Fatal(event, err)
		}
		if _, err = parseDelivery([]byte(body+" "), headers, secret); err == nil {
			t.Fatal("tampering accepted")
		}
	}
}
