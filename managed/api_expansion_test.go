package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"net/http"
	"reflect"
	"testing"
)

func TestManagedExpansionAPIContracts(t *testing.T) {
	t.Run("cancel", func(t *testing.T) {
		c := contractClient(t, "SessionService", "Cancel")
		res, err := c.Sessions.Cancel(context.Background(), "segment /?%#", managed.SessionCancelParams{})
		if err != nil {
			t.Fatal(err)
		}
		checkFields(t, reflect.ValueOf(res), "Session.Cancel")
	})
	t.Run("list_runs", func(t *testing.T) {
		c := contractClient(t, "DeploymentScopedRunService", "List")
		res, err := c.Deployments.Runs.List(context.Background(), "segment /?%#", managed.DeploymentScopedRunListParams{})
		if err != nil {
			t.Fatal(err)
		}
		checkFields(t, reflect.ValueOf(res), "Deployment.Runs.List")
	})
	t.Run("get_run", func(t *testing.T) {
		c := contractClient(t, "DeploymentScopedRunService", "Get")
		res, err := c.Deployments.Runs.Get(context.Background(), "segment /?%#", "segment /?%#", managed.DeploymentScopedRunGetParams{})
		if err != nil {
			t.Fatal(err)
		}
		checkFields(t, reflect.ValueOf(res), "Deployment.Runs.Get")
	})
}

func TestManagedSessionCancelAcknowledgement(t *testing.T) {
	for _, status := range []int{200, 202} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			c := testClient(func(req *http.Request) (*http.Response, error) {
				if req.Method != "POST" || req.URL.Path != "/api/v1/cloud/sessions/sess_one/cancel" {
					t.Fatal(req.URL)
				}
				return reply(req, status, `{"id":"sess_one","type":"session","status":"canceling"}`), nil
			})
			ack, err := c.Sessions.Cancel(context.Background(), "sess_one", managed.SessionCancelParams{})
			if err != nil {
				t.Fatal(err)
			}
			if ack.ID != "sess_one" || ack.Type != "session" || ack.Status != "canceling" {
				t.Fatalf("ack=%+v", ack)
			}
		})
	}
}

func TestManagedDeploymentScopedRunPagination(t *testing.T) {
	calls := 0
	c := testClient(func(req *http.Request) (*http.Response, error) {
		calls++
		if req.URL.Path != "/api/v1/cloud/deployments/dep_one/runs" || req.URL.Query().Get("triggered_after") != "2026-06-01T00:00:00Z" || req.Header.Get("Qoder-Workspace-Id") != "workspace_one" {
			t.Fatalf("request=%v %v", req.URL, req.Header)
		}
		if calls == 2 && req.URL.Query().Get("page") != "opaque +/=?" {
			t.Fatal(req.URL)
		}
		if calls == 1 {
			return reply(req, 200, `{"data":[],"has_more":true,"next_page":"opaque +/=?"}`), nil
		}
		return reply(req, 200, `{"data":[{"id":"drun_one","type":"deployment_run","deployment_id":"dep_one"}],"has_more":false,"next_page":null}`), nil
	})
	pager := c.Deployments.Runs.ListAutoPaging(context.Background(), "dep_one", managed.DeploymentScopedRunListParams{TriggeredAfter: managed.String("2026-06-01T00:00:00Z"), WorkspaceID: managed.String("workspace_one")})
	if !pager.Next() || pager.Current().ID != "drun_one" || pager.Next() || pager.Err() != nil || calls != 2 {
		t.Fatalf("calls=%d err=%v", calls, pager.Err())
	}
}
