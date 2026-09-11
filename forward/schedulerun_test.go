package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestScheduleRunAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listScheduleRuns", status)
				res, err := client.ScheduleRuns.List(context.Background(), contractParams[forward.ScheduleRunListParams](t, "listScheduleRuns"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Get", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "getScheduleRun", status)
				res, err := client.ScheduleRuns.Get(context.Background(), pathSegment, contractParams[forward.ScheduleRunGetParams](t, "getScheduleRun"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
}
