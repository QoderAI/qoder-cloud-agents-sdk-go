package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestScheduleAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listSchedules", status)
				res, err := client.Schedules.List(context.Background(), contractParams[forward.ScheduleListParams](t, "listSchedules"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("New", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "createSchedule", status)
				res, err := client.Schedules.New(context.Background(), contractParams[forward.ScheduleNewParams](t, "createSchedule"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("ArchiveMany", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "archiveSchedules", status)
				res, err := client.Schedules.ArchiveMany(context.Background(), forward.ScheduleArchiveManyParams{
					Scope:          "by_schedule_ids",
					ScheduleIDs:    []string{"sdk-contract"},
					IdempotencyKey: forward.String("sdk-contract"),
				})
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
				client := contractClient(t, "getSchedule", status)
				res, err := client.Schedules.Get(context.Background(), pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Update", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "updateSchedule", status)
				res, err := client.Schedules.Update(context.Background(), pathSegment, contractParams[forward.ScheduleUpdateParams](t, "updateSchedule"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Archive", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "archiveSchedule", status)
				res, err := client.Schedules.Archive(context.Background(), pathSegment, contractParams[forward.ScheduleArchiveParams](t, "archiveSchedule"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Pause", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "pauseSchedule", status)
				res, err := client.Schedules.Pause(context.Background(), pathSegment, contractParams[forward.SchedulePauseParams](t, "pauseSchedule"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Run", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "runSchedule", status)
				res, err := client.Schedules.Run(context.Background(), pathSegment, contractParams[forward.ScheduleRunParams](t, "runSchedule"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Unpause", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "unpauseSchedule", status)
				res, err := client.Schedules.Unpause(context.Background(), pathSegment, contractParams[forward.ScheduleUnpauseParams](t, "unpauseSchedule"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
}
