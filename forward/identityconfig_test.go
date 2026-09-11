package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestIdentityConfigAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listIdentityConfigs", status)
				res, err := client.Identities.Configs.List(context.Background(), pathSegment, contractParams[forward.IdentityConfigListParams](t, "listIdentityConfigs"))
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
				client := contractClient(t, "getIdentityConfig", status)
				res, err := client.Identities.Configs.Get(context.Background(), pathSegment, pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Upsert", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "upsertIdentityConfig", status)
				res, err := client.Identities.Configs.Upsert(context.Background(), pathSegment, pathSegment, contractParams[forward.IdentityConfigUpsertParams](t, "upsertIdentityConfig"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("GetEffective", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "getEffectiveConfig", status)
				res, err := client.Identities.Configs.GetEffective(context.Background(), pathSegment, pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
}
