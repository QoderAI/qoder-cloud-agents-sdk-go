package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
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
		Name:        "custom-tools",
		Description: "让 Agent 调用订单查询工具，在本地 Go 函数中执行，再将结果回传给 Agent。",
		Run:         func(ctx context.Context, r *live.Run) error { return run(ctx, r, s) },
	}))
}

type order struct {
	ID             string `json:"order_id"`
	Status         string `json:"status"`
	TrackingNumber string `json:"tracking_number"`
}

func run(ctx context.Context, r *live.Run, s *managedutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	if err := s.SelectModel(ctx, r); err != nil {
		return err
	}
	r.Step("创建带 lookup_order 自定义工具的 Agent")
	agent, err := s.Client.Agents.New(ctx, managed.AgentNewParams{
		Name:   live.Name("custom-tools"),
		Model:  managed.ManagedAgentsModelConfigParams{ID: s.Model},
		System: managed.String("查询订单时必须调用 lookup_order。拿到工具返回后，用中文回答订单状态和完整运单号，不得编造。"),
		Tools: []managed.AgentNewParamsToolUnion{{OfCustom: &managed.ManagedAgentsCustomToolParams{
			Type: "custom", Name: "lookup_order", Description: "根据订单 ID 查询订单状态和运单号。",
			InputSchema: managed.ManagedAgentsCustomToolInputSchemaParam{
				Properties: map[string]any{"order_id": map[string]any{"type": "string", "description": "待查询的订单 ID"}},
				Required:   []string{"order_id"},
			},
		}}},
	})
	if err != nil {
		return err
	}
	r.Track("agent", agent.ID, func(ctx context.Context) error {
		_, err := s.Client.Agents.Archive(ctx, agent.ID, managed.AgentArchiveParams{})
		return err
	})
	session, err := s.NewSession(ctx, r, managed.SessionNewParams{EnvironmentID: env, Agent: managed.SessionNewParamsAgentUnion{OfString: managed.String(agent.ID)}})
	if err != nil {
		return err
	}
	// The tracking number exists only in this Go process, not in the model prompt.
	data := order{ID: "order-" + live.Marker(), Status: "已发货", TrackingNumber: "track-" + live.Marker()}
	prompt := "请查询订单 " + data.ID + "，告诉我订单状态和完整运单号。"
	r.Step("发送订单查询消息")
	r.Message("user", prompt)
	after, err := s.Send(ctx, session, prompt)
	if err != nil {
		return err
	}
	return runTools(ctx, r, s, session, after, data)
}

// lookupOrder is the business-code boundary. Replace the in-memory lookup with
// your own service or database query; only explicitly registered tools execute.
func lookupOrder(call managed.ManagedAgentsAgentCustomToolUseEvent, data order) (string, bool) {
	if call.Name != "lookup_order" {
		return "unknown tool; use lookup_order", true
	}
	id, ok := call.Input["order_id"].(string)
	if !ok || id == "" {
		return "order_id must be a non-empty string", true
	}
	if id != data.ID {
		return "order not found", true
	}
	body, _ := json.Marshal(data)
	return string(body), false
}

type toolTurn struct {
	result live.TurnResult
	calls  []managed.ManagedAgentsAgentCustomToolUseEvent
	cursor string
}

func readToolTurn(ctx context.Context, r *live.Run, s *managedutil.Example, session, after string) (toolTurn, error) {
	turn := toolTurn{result: live.TurnResult{LastID: after}, cursor: after}
	stream := s.Client.Sessions.Events.StreamEvents(ctx, session, managed.SessionEventStreamParams{}, option.WithHeader("Last-Event-ID", after), option.WithRequestTimeout(s.Config.Timeout))
	defer stream.Close()
	pending := map[string]managed.ManagedAgentsAgentCustomToolUseEvent{}
	for stream.Next() {
		event := stream.Current()
		if event.ID != "" {
			turn.cursor = event.ID
		}
		if event.Type == "agent.custom_tool_use" {
			var call managed.ManagedAgentsAgentCustomToolUseEvent
			if err := json.Unmarshal([]byte(event.RawJSON()), &call); err != nil {
				return turn, err
			}
			if call.ID == "" || call.SessionThreadID != "" {
				return turn, errors.New("this example expects custom tool calls on the primary session")
			}
			pending[call.ID] = call
			r.Log("tool", call.Name)
			continue
		}
		// requires_action is a pause for tool results, not a completed answer.
		if event.Type == "session.status_idle" && event.StopReason.Type == "requires_action" {
			if event.ID == "" || len(event.StopReason.EventIDs) == 0 || len(event.StopReason.EventIDs) > 8 {
				return turn, errors.New("invalid or excessive pending tool calls")
			}
			seen := map[string]bool{}
			for _, id := range event.StopReason.EventIDs {
				call, ok := pending[id]
				if !ok {
					return turn, fmt.Errorf("requires_action references unknown custom tool call %s", id)
				}
				if !seen[id] {
					turn.calls = append(turn.calls, call)
					seen[id] = true
				}
			}
			return turn, nil
		}
		if err := r.ObserveTurn(&turn.result, event.RawJSON()); err != nil {
			return turn, err
		}
		if turn.result.Complete {
			return turn, nil
		}
	}
	if err := stream.Err(); err != nil {
		return turn, err
	}
	return turn, errors.New("stream ended before a final answer or tool-result request")
}

func runTools(ctx context.Context, r *live.Run, s *managedutil.Example, session, after string, data order) error {
	stopWaiting := r.Wait("等待 Agent 调用本地工具并给出最终回答")
	defer stopWaiting()
	successfulCalls := 0
	for round := 0; round < 8; round++ {
		turn, err := readToolTurn(ctx, r, s, session, after)
		if err != nil {
			return err
		}
		if turn.result.Complete {
			if successfulCalls == 0 {
				return errors.New("agent completed without a successful local tool call")
			}
			if err := turn.result.Verify([]string{data.Status, data.TrackingNumber}, false); err != nil {
				return err
			}
			r.Log("checked", "最终回答包含仅保存在 Go 工具中的订单状态和运单号。")
			return nil
		}
		results := make([]managed.ManagedAgentsEventParamsUnion, 0, len(turn.calls))
		for _, call := range turn.calls {
			text, isError := lookupOrder(call, data)
			if !isError {
				successfulCalls++
			}
			r.Log("info", fmt.Sprintf("本地执行 %s，回传结果（is_error=%t）", call.Name, isError))
			results = append(results, managed.ManagedAgentsEventParamsUnion{OfUserCustomToolResult: &managed.ManagedAgentsUserCustomToolResultEventParams{
				Type: "user.custom_tool_result", CustomToolUseID: call.ID, IsError: managed.Bool(isError),
				Content: []managed.ManagedAgentsUserCustomToolResultEventParamsContentUnion{{OfText: &managed.ManagedAgentsTextBlockParam{Type: "text", Text: text}}},
			}})
		}
		response, err := s.Client.Sessions.Events.Send(ctx, session, managed.SessionEventSendParams{Events: results}, option.WithHeader("Idempotency-Key", live.Name("tool-results")))
		if err != nil {
			return err
		}
		if len(response.Data) != len(results) {
			return errors.New("server did not acknowledge every tool result")
		}
		// Resume after the idle event so a fast final reply cannot be skipped.
		after = turn.cursor
	}
	return errors.New("custom tool loop exceeded eight rounds")
}
