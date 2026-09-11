package main

import (
	"context"
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
		Name:        "session",
		Description: "创建 Agent 和 Session，发送一条消息，通过 SSE 接收助手回复。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

// run demonstrates the Agent -> Session lifecycle and SSE completion.
func run(ctx context.Context, r *live.Run, s *managedutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	agent, err := s.Agent(ctx, r, managed.AgentNewParams{})
	if err != nil {
		return err
	}
	session, err := s.NewSession(ctx, r, managed.SessionNewParams{EnvironmentID: env, Agent: managed.SessionNewParamsAgentUnion{OfString: managed.String(agent)}})
	if err != nil {
		return err
	}
	marker := live.Marker()
	return s.Turn(ctx, r, session, "请用一句话介绍你能提供什么帮助，并在回复末尾原样附上："+marker, []string{marker}, false, true)
}
