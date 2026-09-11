package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestVaultAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listVault", status)
				res, err := client.Vaults.List(context.Background(), contractParams[forward.VaultListParams](t, "listVault"))
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
				client := contractClient(t, "createVault", status)
				res, err := client.Vaults.New(context.Background(), contractParams[forward.VaultNewParams](t, "createVault"))
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
				client := contractClient(t, "getVault", status)
				res, err := client.Vaults.Get(context.Background(), pathSegment)
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
				client := contractClient(t, "deleteVault", status)
				checkError(t, status, client.Vaults.Delete(context.Background(), pathSegment))
			})
		}
	})
}
