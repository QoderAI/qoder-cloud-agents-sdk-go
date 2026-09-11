package managed_test

import (
	"context"
	"encoding/json"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestEnvironmentWorkGet(t *testing.T) {
	client := contractClient(t, "EnvironmentWorkService", "Get")
	response, err := client.Environments.Work.Get(context.Background(), "segment /?%#", managed.EnvironmentWorkGetParams{EnvironmentID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "EnvironmentWork.Get")
}

func TestEnvironmentWorkUpdate(t *testing.T) {
	client := contractClient(t, "EnvironmentWorkService", "Update")
	response, err := client.Environments.Work.Update(context.Background(), "segment /?%#", managed.EnvironmentWorkUpdateParams{EnvironmentID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "EnvironmentWork.Update")
}

func TestEnvironmentWorkList(t *testing.T) {
	client := contractClient(t, "EnvironmentWorkService", "List")
	response, err := client.Environments.Work.List(context.Background(), "segment /?%#", managed.EnvironmentWorkListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "EnvironmentWork.List")
}

func TestEnvironmentWorkAck(t *testing.T) {
	client := contractClient(t, "EnvironmentWorkService", "Ack")
	response, err := client.Environments.Work.Ack(context.Background(), "segment /?%#", managed.EnvironmentWorkAckParams{EnvironmentID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "EnvironmentWork.Ack")
}

func TestEnvironmentWorkHeartbeat(t *testing.T) {
	client := contractClient(t, "EnvironmentWorkService", "Heartbeat")
	response, err := client.Environments.Work.Heartbeat(context.Background(), "segment /?%#", managed.EnvironmentWorkHeartbeatParams{EnvironmentID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "EnvironmentWork.Heartbeat")
}

func TestEnvironmentWorkPoll(t *testing.T) {
	client := contractClient(t, "EnvironmentWorkService", "Poll")
	response, err := client.Environments.Work.Poll(context.Background(), "segment /?%#", managed.EnvironmentWorkPollParams{})
	if err != nil {
		t.Fatal(err)
	}
	if response != nil {
		t.Fatal("empty work queue must return nil work")
	}
}

func TestEnvironmentWorkStats(t *testing.T) {
	client := contractClient(t, "EnvironmentWorkService", "Stats")
	response, err := client.Environments.Work.Stats(context.Background(), "segment /?%#", managed.EnvironmentWorkStatsParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "EnvironmentWork.Stats")
}

func TestEnvironmentWorkStop(t *testing.T) {
	client := contractClient(t, "EnvironmentWorkService", "Stop")
	response, err := client.Environments.Work.Stop(context.Background(), "segment /?%#", managed.EnvironmentWorkStopParams{EnvironmentID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "EnvironmentWork.Stop")
}

func TestQoderNamesAndWorkerIdentity(t *testing.T) {
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/v1/cloud/environments/env/work/poll" {
			t.Fatal(r.URL.Path)
		}
		if r.Header.Get("Worker-ID") != "worker-01" {
			t.Fatal("worker identity not sent using the documented header")
		}
		if r.Header.Get("X-Qoder-Beta") != "preview" {
			t.Fatal("Qoder preview header lost")
		}
		for header := range r.Header {
			if strings.Contains(strings.ToLower(header), "worker") && header != "Worker-Id" {
				t.Fatalf("unexpected worker header: %s", header)
			}
		}
		return reply(r, 200, `{"data":[]}`), nil
	})
	_, err := c.Environments.Work.Poll(context.Background(), "env", managed.EnvironmentWorkPollParams{QoderWorkerID: managed.String("worker-01"), Betas: []managed.QoderBeta{"preview"}})
	if err != nil {
		t.Fatal(err)
	}
	skill := managed.ManagedAgentsQoderSkillParams{SkillID: "skill", Type: "qoder"}
	body := jsonObject(t, skill)
	if body["type"] != "qoder" || body["skill_id"] != "skill" {
		t.Fatal(body)
	}
	var union managed.ManagedAgentsAgentSkillUnion
	if err := json.Unmarshal([]byte(`{"type":"qoder","skill_id":"skill","version":"1"}`), &union); err != nil {
		t.Fatal(err)
	}
	if union.AsQoder().SkillID != "skill" {
		t.Fatal("Qoder skill union variant was lost")
	}
	if _, ok := union.AsAny().(managed.ManagedAgentsQoderSkill); !ok {
		t.Fatal("Qoder skill discriminator was lost")
	}
}
