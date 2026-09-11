package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestSkillAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listSkill", status)
				res, err := client.Skills.List(context.Background(), contractParams[forward.SkillListParams](t, "listSkill"))
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
				client := contractClient(t, "createSkill", status)
				res, err := client.Skills.New(context.Background(), contractParams[forward.SkillNewParams](t, "createSkill"))
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
				client := contractClient(t, "getSkill", status)
				res, err := client.Skills.Get(context.Background(), pathSegment, contractParams[forward.SkillGetParams](t, "getSkill"))
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
				client := contractClient(t, "updateSkill", status)
				res, err := client.Skills.Update(context.Background(), pathSegment, contractParams[forward.SkillUpdateParams](t, "updateSkill"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Delete", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "deleteSkill", status)
				checkError(t, status, client.Skills.Delete(context.Background(), pathSegment))
			})
		}
	})
}
