package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func main() {
	config, err := live.Load("managed")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	client := managed.NewClient(config.Options()...)
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "models",
		Description: "查询当前账号可用的模型。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, client, config)
		},
	}))
}

func run(ctx context.Context, r *live.Run, client managed.Client, config live.Config) error {
	r.Step("查询账号可用的模型")
	models, err := client.Models.List(ctx, managed.ModelListParams{})
	if err != nil {
		return err
	}
	var enabled []string
	for _, model := range models.Data {
		if model.IsEnabled {
			enabled = append(enabled, model.ID)
		}
	}
	model, err := live.ChooseModel(config.Model, enabled)
	if err != nil {
		return err
	}
	r.Log("info", fmt.Sprintf("已启用 %d 个模型，本次使用：%s", len(enabled), model))
	r.Log("info", "模型列表："+strings.Join(enabled, ", "))
	return nil
}
