package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/managedutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func main() {
	config, err := live.Load("managed")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	s := &managedutil.Example{Client: managed.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "deployment",
		Description: "创建 Deployment，手动触发一次运行并检查助手回复。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

// deployment runs once without a recurring schedule and validates actual model output.
func run(ctx context.Context, r *live.Run, s *managedutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	agent, err := s.Agent(ctx, r, managed.AgentNewParams{})
	if err != nil {
		return err
	}
	marker := live.Marker()
	r.Log("info", "本次运行将发送以下初始消息：")
	r.Message("user", "Reply with exactly "+marker)
	r.Step("创建 Deployment，设置运行时的初始消息")
	deployment, err := s.Client.Deployments.New(ctx, managed.DeploymentNewParams{Name: live.Name("deployment"), EnvironmentID: env, Agent: managed.DeploymentNewParamsAgentUnion{OfString: managed.String(agent)}, InitialEvents: []managed.ManagedAgentsDeploymentInitialEventParamsUnion{{OfUserMessage: &managed.ManagedAgentsUserMessageEventParams{Type: "user.message", Content: []managed.ManagedAgentsUserMessageEventParamsContentUnion{{OfText: &managed.ManagedAgentsTextBlockParam{Type: "text", Text: "Reply with exactly " + marker}}}}}}})
	if err != nil {
		return err
	}
	r.Track("deployment", deployment.ID, func(ctx context.Context) error {
		_, err := s.Client.Deployments.Archive(ctx, deployment.ID, managed.DeploymentArchiveParams{})
		return err
	})
	r.Step("手动触发一次 Deployment 运行")
	run, err := s.Client.Deployments.Run(ctx, deployment.ID, managed.DeploymentRunParams{})
	if err != nil {
		return err
	}
	r.Log("info", "本次运行 ID："+run.ID)
	if run.SessionID == "" {
		return errors.New("deployment run returned no session")
	}
	r.Track("session", run.SessionID, func(ctx context.Context) error { return s.FinishSession(ctx, run.SessionID) })
	got, err := s.Client.DeploymentRuns.Get(ctx, run.ID, managed.DeploymentRunGetParams{})
	if err != nil {
		return err
	}
	if got.SessionID != run.SessionID {
		return errors.New("deployment run session ID changed")
	}
	return s.WaitTurn(ctx, r, run.SessionID, "", []string{marker}, false, false)
}
