package forward

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ErrWebhookVerification indicates missing, invalid, expired or mismatched delivery authentication.
var ErrWebhookVerification = errors.New("webhook signature verification failed")

// VerifyWebhookSignature authenticates the original request bytes with the endpoint's
// whsec_ secret. It accepts any matching v1 signature within 300 seconds of now.
// Verify before parsing JSON; deduplicate accepted events using Webhook-ID.
func VerifyWebhookSignature(payload []byte, headers http.Header, secret string) error {
	return verifyWebhookSignatureAt(payload, headers, secret, time.Now().Unix())
}
func verifyWebhookSignatureAt(payload []byte, headers http.Header, secret string, now int64) error {
	get := func(name string) string {
		var found []string
		for k, v := range headers {
			if strings.EqualFold(k, name) {
				found = append(found, v...)
			}
		}
		if len(found) != 1 {
			return ""
		}
		return found[0]
	}
	id, timestamp, signatures := get("Webhook-ID"), get("Webhook-Timestamp"), get("Webhook-Signature")
	if id == "" || strings.ContainsAny(id, ".\r\n") || timestamp == "" || len(timestamp) > 19 || !strings.HasPrefix(secret, "whsec_") {
		return ErrWebhookVerification
	}
	for _, c := range timestamp {
		if c < '0' || c > '9' {
			return ErrWebhookVerification
		}
	}
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || ts < now-300 || ts > now+300 {
		return ErrWebhookVerification
	}
	encoded := strings.TrimPrefix(secret, "whsec_")
	key, err := base64.StdEncoding.Strict().DecodeString(encoded)
	if err != nil || len(key) < 24 || len(key) > 64 || base64.StdEncoding.EncodeToString(key) != encoded {
		return ErrWebhookVerification
	}
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(id + "." + timestamp + "."))
	mac.Write(payload)
	expected := mac.Sum(nil)
	for _, signature := range strings.Fields(signatures) {
		if !strings.HasPrefix(signature, "v1,") {
			continue
		}
		candidate, err := base64.StdEncoding.Strict().DecodeString(strings.TrimPrefix(signature, "v1,"))
		if err == nil && hmac.Equal(candidate, expected) {
			return nil
		}
	}
	return ErrWebhookVerification
}
