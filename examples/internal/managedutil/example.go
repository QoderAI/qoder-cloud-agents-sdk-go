// Package managedutil provides setup and session helpers for the managed examples.
package managedutil

import (
	"context"
	"fmt"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

type Example struct {
	Client managed.Client
	Config live.Config
	Model  string
}

func (s *Example) SelectModel(ctx context.Context, r *live.Run) error {
	r.Step("查询账号可用的模型")
	models, err := s.Client.Models.List(ctx, managed.ModelListParams{})
	if err != nil {
		return err
	}
	var enabled []string
	for _, model := range models.Data {
		if model.IsEnabled {
			enabled = append(enabled, model.ID)
		}
	}
	s.Model, err = live.ChooseModel(s.Config.Model, enabled)
	if err != nil {
		return err
	}
	r.Log("info", fmt.Sprintf("已启用 %d 个模型，本次使用：%s", len(enabled), s.Model))
	return nil
}
