//go:build live

package managed_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestManagedSessionCleanupOffline(t *testing.T) {
	var requests []string
	reads := 0
	s := &managedScenarioSuite{client: testClient(func(r *http.Request) (*http.Response, error) {
		key := r.Method + " " + r.URL.Path
		requests = append(requests, key)
		switch key {
		case "GET /api/v1/cloud/sessions/session":
			reads++
			status := "running"
			if reads > 1 {
				status = "idle"
			}
			return reply(r, 200, `{"id":"session","status":"`+status+`"}`), nil
		case "POST /api/v1/cloud/sessions/session/events":
			var body struct{ Events []struct{ Type string } }
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if len(body.Events) != 1 || body.Events[0].Type != "user.interrupt" {
				t.Fatal("cleanup did not interrupt")
			}
			return reply(r, 200, `{"data":[]}`), nil
		case "DELETE /api/v1/cloud/sessions/session":
			return reply(r, 200, `{"id":"session","type":"session_deleted"}`), nil
		default:
			t.Fatalf("unexpected request %s", key)
			return nil, nil
		}
	})}
	if err := s.finishSession(context.Background(), "session"); err != nil {
		t.Fatal(err)
	}
	want := []string{"GET /api/v1/cloud/sessions/session", "POST /api/v1/cloud/sessions/session/events", "GET /api/v1/cloud/sessions/session", "DELETE /api/v1/cloud/sessions/session"}
	if !reflect.DeepEqual(requests, want) {
		t.Fatalf("cleanup order=%v", requests)
	}
}
