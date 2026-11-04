package forward_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestTemplateAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listTemplates", status)
				res, err := client.Templates.List(context.Background(), contractParams[forward.TemplateListParams](t, "listTemplates"))
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
				client := contractClient(t, "createTemplate", status)
				res, err := client.Templates.New(context.Background(), contractParams[forward.TemplateNewParams](t, "createTemplate"))
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
				client := contractClient(t, "getTemplate", status)
				res, err := client.Templates.Get(context.Background(), pathSegment)
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
				client := contractClient(t, "updateTemplate", status)
				res, err := client.Templates.Update(context.Background(), pathSegment, contractParams[forward.TemplateUpdateParams](t, "updateTemplate"))
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
				client := contractClient(t, "archiveTemplate", status)
				res, err := client.Templates.Archive(context.Background(), pathSegment, contractParams[forward.TemplateArchiveParams](t, "archiveTemplate"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Clone", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "cloneTemplate", status)
				res, err := client.Templates.Clone(context.Background(), pathSegment, contractParams[forward.TemplateCloneParams](t, "cloneTemplate"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
}

func TestTemplateTunnelBinding(t *testing.T) {
	for _, tc := range []struct {
		name   string
		params any
		want   string
	}{
		{"new", forward.TemplateNewParams{TunnelID: forward.String("tnl_one")}, `"tnl_one"`},
		{"keep", forward.TemplateUpdateParams{}, ""},
		{"replace", forward.TemplateUpdateParams{TunnelID: forward.String("tnl_two")}, `"tnl_two"`},
		{"clear", forward.TemplateUpdateParams{TunnelID: param.Null[string]()}, "null"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := json.Marshal(tc.params)
			if err != nil {
				t.Fatal(err)
			}
			var body map[string]json.RawMessage
			if err := json.Unmarshal(raw, &body); err != nil {
				t.Fatal(err)
			}
			if got := string(body["tunnel_id"]); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
	var result forward.Template
	if err := json.Unmarshal([]byte(`{"id":"tmpl_one","tunnel_id":"tnl_one"}`), &result); err != nil {
		t.Fatal(err)
	}
	if result.TunnelID != "tnl_one" || !result.JSON.TunnelID.Valid() {
		t.Fatalf("tunnel_id=%q", result.TunnelID)
	}
}
